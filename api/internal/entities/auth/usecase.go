package auth

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
	"github.com/kevinburke/twilio-go"
	"github.com/stripe/stripe-go/v72"
)

type Usecase interface {
	AuthUser(ctx context.Context, tx pgx.Tx, userUUID string, userID int) (*response.AuthSuccess, error)
	VerifySMS(ctx context.Context, req request.VerifySMS) (*response.AuthConfirm, error)
	SignOut(ctx context.Context, req request.SignOut) error
	SignInSMS(ctx context.Context, req request.SignInSMS) (*response.AuthSuccess, error)
	SignupSMS(ctx context.Context, req request.SignupSMS) (*response.AuthSuccess, error)
	SignupSMSGoat(ctx context.Context, formattedNumber string, req request.SignupSMSGoat) (*response.AuthSuccess, error)
	RefreshToken(ctx context.Context, req request.RefreshToken) (*response.AuthSuccess, error)
	VerifyTwilioCode(ctx context.Context, phone string, code string) (*twilio.CheckPhoneNumber, error)

	CheckPhone(ctx context.Context, req request.Phone) (*response.Check, error)
	CheckUsername(ctx context.Context, req request.Username) (*response.Check, error)
	CheckEmail(ctx context.Context, req request.Email) (*response.Check, error)
	CheckGoatInviteCode(ctx context.Context, req request.GoatInviteCode) (*response.Check, error)

	GetAuthToken(claims BaseClaims, ttlHours int) (string, int64, error)
	GetRefreshToken(claims BaseClaims, ttlHours int) (string, int64, error)

	GenerateGoatInviteCode(ctx context.Context, insertCode bool) (*response.GenerateGoatInviteCode, error)
	GenerateUserInviteCodes(ctx context.Context, exec utils.Executable, userID int) error

	CreateProviderAccount(ctx context.Context, req request.CreateProviderAccount) (*response.CreateProviderAccount, error)
	CreateSendBirdUser(currUser response.User) (string, error)
	CreateStripeCustomer(stripeParams *stripe.CustomerParams) (string, error)

	AddUsersToContacts(ctx context.Context, exec utils.Executable, user response.User) error
}
