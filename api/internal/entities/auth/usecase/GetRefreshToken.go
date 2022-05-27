package service

import (
	"time"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	auth "github.com/VamaSingapore/vama-api/internal/entities/auth"
	"github.com/golang-jwt/jwt"
)

func (svc *usecase) GetRefreshToken(c auth.BaseClaims, ttlHours int) (string, int64, error) {
	// sec * minutes * hours
	ttlTimeResult := time.Duration(60 * 60 * ttlHours)
	ttl := time.Now().Add(time.Second * ttlTimeResult).Unix()
	claims := auth.Claims{
		Data: c,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: ttl,
		},
	}
	refreshTokenKey := []byte(appconfig.Config.Auth.RefreshTokenKey)
	return svc.genHS256JWT(claims, ttl, refreshTokenKey)
}
