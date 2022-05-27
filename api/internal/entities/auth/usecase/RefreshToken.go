package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	auth "github.com/VamaSingapore/vama-api/internal/entities/auth"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"golang.org/x/crypto/scrypt"
)

const refreshErr = "Something went wrong when trying to authenticate you.  Please try signing in again."

func (svc *usecase) RefreshToken(ctx context.Context, req request.RefreshToken) (*response.AuthSuccess, error) {
	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			refreshErr,
			fmt.Sprintf("Could not start transaction. Error: %v", txErr),
		)
	}

	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	tokenClaims, verifyIDTokenErr := svc.token.VerifyIDToken(ctx, tx, req.RefreshToken, constants.REFRESH_TOKEN_TYPE)
	if verifyIDTokenErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			401,
			http.StatusUnauthorized,
			refreshErr,
			fmt.Sprintf("Failed to validate refresh token: %v", verifyIDTokenErr),
		)
	}
	if tokenClaims.ExpiresAt < time.Now().Unix() {
		return nil, httperr.NewCtx(
			ctx,
			401,
			http.StatusUnauthorized,
			refreshErr,
			"The refresh token has expired.",
		)
	}

	costParameter := 16384
	r := 8
	p := 1
	keyLen := 32
	refreshTokenHash, refreshHashErr := scrypt.Key([]byte(req.RefreshToken), []byte(tokenClaims.Data.UUID), costParameter, r, p, keyLen)
	if refreshHashErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			refreshErr,
			fmt.Sprintf("Error hashing refresh token. Err: %v", refreshHashErr),
		)
	}

	rTokenNotRevoked := new(bool)
	getFromCacheErr := svc.repo.GetFromCache(ctx, tx, constants.REFRESH_TOKEN_TYPE, fmt.Sprintf("%x", refreshTokenHash), &rTokenNotRevoked)
	if getFromCacheErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			refreshErr,
			fmt.Sprintf("Error getting refresh token from cache: %v", getFromCacheErr),
		)
	}

	if rTokenNotRevoked == nil || !*rTokenNotRevoked {
		return nil, httperr.NewCtx(
			ctx,
			401,
			http.StatusUnauthorized,
			refreshErr,
			"Refresh token has been revoked",
		)
	}

	c := auth.BaseClaims{
		UUID: tokenClaims.Data.UUID,
	}

	newToken, ttlToken, getAuthTokenErr := svc.GetAuthToken(c, ttlTokenHours)
	if getAuthTokenErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			refreshErr,
			fmt.Sprintf("Error getting auth token. Err: %v", getAuthTokenErr),
		)
	}

	newTokenHash, hashAuthTokenErr := scrypt.Key([]byte(newToken), []byte(tokenClaims.Data.UUID), costParameter, r, p, keyLen)
	if hashAuthTokenErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			refreshErr,
			fmt.Sprintf("Error hashing auth token. Err: %v", hashAuthTokenErr),
		)
	}

	user, getUserErr := svc.user.GetUserByUUID(ctx, tx, c.UUID)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			refreshErr,
			fmt.Sprintf("Error getting user by uuid: %v", getUserErr),
		)
	} else if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			refreshErr,
			fmt.Sprintf("User %s was nil in RefreshToken", c.UUID),
		)
	}

	accessCacheErr := svc.repo.SaveToCacheForHours(ctx, tx, constants.ACCESS_TOKEN_TYPE, fmt.Sprintf("%x", newTokenHash), ttlTokenHours, true, user.ID)
	if accessCacheErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			refreshErr,
			fmt.Sprintf("Error getting access token from cache. Err: %v", accessCacheErr),
		)
	}

	resp := response.AuthSuccess{}
	resp.Credentials.AccessToken = newToken
	resp.Credentials.AccessTokenExpiresAt = ttlToken

	commit = true
	return &resp, nil
}
