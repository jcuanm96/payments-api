package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	auth "github.com/VamaSingapore/vama-api/internal/entities/auth"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/scrypt"
)

func (svc *usecase) AuthUser(ctx context.Context, tx pgx.Tx, userUUID string, userID int) (*response.AuthSuccess, error) {
	if userUUID == "" {
		return nil, errors.New("user has an invalid uuid")
	}
	c := auth.BaseClaims{
		UUID: userUUID,
	}

	token, ttlToken, err := svc.GetAuthToken(c, ttlTokenHours)
	if err != nil {
		return nil, err
	}

	// Hashing the access and refresh tokens and storing them in cache
	costParameter := 16384
	r := 8
	p := 1
	keyLen := 32
	tokenHash, hashErr := scrypt.Key([]byte(token), []byte(userUUID), costParameter, r, p, keyLen)
	if hashErr != nil {
		return nil, hashErr
	}

	accessCacheErr := svc.repo.SaveToCacheForHours(ctx, tx, constants.ACCESS_TOKEN_TYPE, fmt.Sprintf("%x", tokenHash), ttlTokenHours, true, userID)
	if accessCacheErr != nil {
		return nil, accessCacheErr
	}

	rToken, ttlRefreshToken, getRefreshTokenErr := svc.GetRefreshToken(c, ttlRefreshTokenHours)
	if getRefreshTokenErr != nil {
		return nil, getRefreshTokenErr
	}
	rTokenHash, hashRefreshErr := scrypt.Key([]byte(rToken), []byte(userUUID), costParameter, r, p, keyLen)
	if hashRefreshErr != nil {
		return nil, hashRefreshErr
	}

	refreshCacheErr := svc.repo.SaveToCacheForHours(ctx, tx, constants.REFRESH_TOKEN_TYPE, fmt.Sprintf("%x", rTokenHash), ttlRefreshTokenHours, true, userID)
	if refreshCacheErr != nil {
		return nil, refreshCacheErr
	}

	resp := response.AuthSuccess{}
	resp.Credentials.AccessToken = token
	resp.Credentials.AccessTokenExpiresAt = ttlToken
	resp.Credentials.RefreshToken = rToken
	resp.Credentials.RefreshTokenExpiresAt = ttlRefreshToken
	return &resp, nil
}
