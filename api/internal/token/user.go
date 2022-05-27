package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/scrypt"
)

const (
	authorizationScheme = "Bearer"
)

type Service interface {
	GetRedisClient() vredis.Client
	MasterNode() *pgxpool.Pool
	VerifyIDToken(ctx context.Context, runnable utils.Runnable, idToken string, tokenType string) (*Claims, error)
	GetTokenFromAuthHeader(headerVal string) (string, error)
}

type svc struct {
	db             *pgxpool.Pool
	rdb            vredis.Client
	accessTokenKey string
}

func NewService(
	db *pgxpool.Pool,
	rdb vredis.Client,
	accessTokenKey string,
) Service {
	return &svc{
		accessTokenKey: accessTokenKey,
		rdb:            rdb,
		db:             db,
	}
}

func (s *svc) GetRedisClient() vredis.Client {
	return s.rdb
}

func (s *svc) MasterNode() *pgxpool.Pool {
	return s.db
}
func (s *svc) VerifyIDToken(ctx context.Context, runnable utils.Runnable, idToken string, tokenType string) (*Claims, error) {
	var idTokenKey string
	if tokenType == constants.ACCESS_TOKEN_TYPE {
		idTokenKey = appconfig.Config.Auth.AccessTokenKey
	} else {
		idTokenKey = appconfig.Config.Auth.RefreshTokenKey
	}
	valid, verifyTokenErr := s.verifyToken(ctx, runnable, idToken, []byte(idTokenKey), tokenType)
	if verifyTokenErr != nil {
		return &Claims{}, verifyTokenErr
	}
	if !valid {
		return nil, errors.New("token invalid")
	}
	parser := new(jwt.Parser)
	c := &Claims{}
	claims := &Claims{}
	token, _, err := parser.ParseUnverified(idToken, claims)
	if err != nil {
		return &Claims{}, err
	}
	err = ConvertViaJSON(token.Claims, &c)
	if err != nil {
		return &Claims{}, err
	}
	return claims, nil
}

func ConvertViaJSON(from, to interface{}) error {
	data, err := json.Marshal(from)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, to)
}

func (s *svc) GetTokenFromAuthHeader(headerVal string) (string, error) {
	if !strings.Contains(headerVal, authorizationScheme) {
		return "", errors.New("missing bearer token")
	}
	headerParts := strings.Split(headerVal, " ")
	if len(headerParts) != 2 {
		return "", errors.New("malformed bearer token")
	}
	return headerParts[1], nil
}

func (s *svc) verifyToken(ctx context.Context, runnable utils.Runnable, tokenString string, key []byte, tokenType string) (bool, error) {
	if tokenString == "" {
		return false, nil
	}
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false, errors.New("token does not have 3 parts")
	}
	err := jwt.SigningMethodHS256.Verify(strings.Join(parts[0:2], "."), parts[2], key)

	if err != nil {
		return false, nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return false, err
	}

	claims := token.Claims.(jwt.MapClaims)
	data := claims["data"].(map[string]interface{})
	uuid := data["uuid"].(string)

	costParameter := 16384
	r := 8
	p := 1
	keyLen := 32
	tokenHash, err := scrypt.Key([]byte(tokenString), []byte(uuid), costParameter, r, p, keyLen)
	if err != nil {
		return false, err
	}

	// Check that the user hasn't already logged out (cleared their tokens)
	var accessTokenRevoked bool
	getFromCacheErr := baserepo.GetFromCache(ctx, runnable, s.GetRedisClient(), tokenType, fmt.Sprintf("%x", tokenHash), &accessTokenRevoked)
	if getFromCacheErr != nil {
		vlog.Errorf(ctx, "Error getting from cache in verifyToken: %v", getFromCacheErr)
		return false, getFromCacheErr
	}
	if !accessTokenRevoked {
		return false, nil
	}

	return token.Valid, nil
}

type Claims struct {
	Data BaseClaims `json:"data"`
	jwt.StandardClaims
}

type BaseClaims struct {
	UUID string `json:"uuid"`
}
