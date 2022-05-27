package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/token"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/codes"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"golang.org/x/crypto/scrypt"
	"google.golang.org/api/idtoken"
)

func MustParseClaims(tokenSvc token.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		tx, txErr := tokenSvc.MasterNode().Begin(ctx)
		if txErr != nil {
			return httperr.New(
				500,
				http.StatusInternalServerError,
				"Something went wrong when trying to authenticate you.").
				SetDetail(txErr).Send(c)
		}
		defer tx.Rollback(ctx)

		idToken, err := tokenSvc.GetTokenFromAuthHeader(c.Get("Authorization"))
		if err != nil {
			return err
		}
		token, err := tokenSvc.VerifyIDToken(ctx, tx, idToken, constants.ACCESS_TOKEN_TYPE)
		if err != nil {
			return httperr.New(codes.Omit, http.StatusUnauthorized, "invalid JWT auth token").SetDetail(err).Send(c)
		}
		if token.ExpiresAt < time.Now().Unix() {
			return httperr.New(codes.Omit, http.StatusUnauthorized, "auth token expired")
		}

		authTkn := strings.TrimPrefix(idToken, "Bearer ")

		costParameter := 16384
		r := 8
		p := 1
		keyLen := 32
		tokenHash, err := scrypt.Key([]byte(authTkn), []byte(token.Data.UUID), costParameter, r, p, keyLen)
		if err != nil {
			return err
		}

		isTokenValid := new(bool)
		getFromCacheErr := baserepo.GetFromCache(ctx, tx, tokenSvc.GetRedisClient(), constants.ACCESS_TOKEN_TYPE, fmt.Sprintf("%x", tokenHash), &isTokenValid)
		if getFromCacheErr != nil {
			vlog.Errorf(ctx, "Error getting from cache in MustParseClaims: %v", getFromCacheErr)
			return httperr.New(codes.Omit, http.StatusUnauthorized, "Something went wrong trying to authenticate you.")
		}
		if isTokenValid == nil || !*isTokenValid {
			return httperr.New(codes.Omit, http.StatusUnauthorized, "auth token revoked")
		}

		c.Locals("uid", token.Data.UUID)

		// Explicitly commit rather than defer due to c.Next()
		// executing more code
		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			vlog.Errorf(ctx, "Error committing tx: %v", commitErr)
		}
		return c.Next()
	}
}

func ParseGCPServiceAuthorization(tokenSvc token.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idToken, err := tokenSvc.GetTokenFromAuthHeader(c.Get("Authorization"))
		if err != nil {
			return err
		}

		payload, validateErr := idtoken.Validate(context.Background(), idToken, appconfig.Config.Gcloud.APIBaseURL)
		if validateErr != nil {
			return httperr.New(
				401,
				http.StatusUnauthorized,
				"Error validating gcp service request",
				fmt.Sprintf("Error validating gcp service request: %v", validateErr),
			)
		} else if payload == nil {
			return httperr.New(
				401,
				http.StatusUnauthorized,
				"Unauthorized gcp service",
				"GCP service validation payload was nil",
			)
		}

		return c.Next()
	}
}

type LimiterConfig struct {
	Max        int
	Expiration time.Duration
}

func Limiter(config *LimiterConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        config.Max,
		Expiration: config.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return httperr.New(codes.Omit, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests)).Send(c)
		},
	})
}
