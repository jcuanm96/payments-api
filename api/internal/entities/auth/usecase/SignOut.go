package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/scrypt"
)

const signoutErr = "Something went wrong when trying to sign out.  Please try again."

func (svc *usecase) SignOut(ctx context.Context, req request.SignOut) error {
	refreshToken, parseErr := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(appconfig.Config.Auth.RefreshTokenKey), nil
	})

	if parseErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signoutErr,
			fmt.Sprintf("Error parsing refresh token. Err: %v", parseErr),
		)
	}

	claims := refreshToken.Claims.(jwt.MapClaims)
	data := claims["data"].(map[string]interface{})
	uuid := data["uuid"].(string)

	costParameter := 16384
	r := 8
	p := 1
	keyLen := 32
	refreshTokenHash, hashErr := scrypt.Key([]byte(req.RefreshToken), []byte(uuid), costParameter, r, p, keyLen)
	if hashErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signoutErr,
			fmt.Sprintf("Error occurred trying to hash refresh token. Err: %v", hashErr),
		)
	}

	clearCachErr := svc.repo.ClearCache(ctx, constants.REFRESH_TOKEN_TYPE, fmt.Sprintf("%x", refreshTokenHash))
	if clearCachErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signoutErr,
			fmt.Sprintf("Error occurred trying to clear cache. Err: %v", clearCachErr),
		)
	}

	return nil
}
