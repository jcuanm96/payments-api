package service

import (
	"time"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	auth "github.com/VamaSingapore/vama-api/internal/entities/auth"
	"github.com/golang-jwt/jwt"
)

func (svc *usecase) GetAuthToken(c auth.BaseClaims, ttlHours int) (string, int64, error) {
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
	authTokenKey := []byte(appconfig.Config.Auth.AccessTokenKey)
	return svc.genHS256JWT(claims, ttl, authTokenKey)
}

func (svc *usecase) genHS256JWT(claims jwt.Claims, ttl int64, signedKey []byte) (string, int64, error) {
	tokenString := ""

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signedKey)
	if err != nil {
		return "", 0, err
	}
	return tokenString, ttl, nil
}
