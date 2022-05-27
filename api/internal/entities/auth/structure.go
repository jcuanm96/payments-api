package auth

import "github.com/golang-jwt/jwt"

type BaseClaims struct {
	UUID string `json:"uuid"`
}
type Claims struct {
	Data BaseClaims `json:"data"`
	jwt.StandardClaims
}
