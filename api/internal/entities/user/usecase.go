package user

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
)

type Usecase interface {
	ListUsers(ctx context.Context, req request.ListUsers) (*response.ListUsers, error)

	GetCurrentUser(c context.Context) (*response.User, error)
	GetMyProfile(ctx context.Context) (*response.MyProfile, error)
	GetUserByID(ctx context.Context, id int) (*response.User, error)
	GetUserByUUID(ctx context.Context, runnable utils.Runnable, uuid string) (*response.User, error)
	GetUserByEmail(ctx context.Context, runnable utils.Runnable, email string) (*response.User, error)
	GetUserByPhone(ctx context.Context, runnable utils.Runnable, phone string) (*response.User, error)
	GetGoatProfile(ctx context.Context, req request.GetGoatProfile) (*response.GetGoatProfile, error)
	SoftDeleteUser(ctx context.Context, userID int) error
	HardDeleteUser(ctx context.Context) error
	UndeleteUser(ctx context.Context, userID int) error
	ReportUser(ctx context.Context, req request.ReportUser) error
	GetUserByUsername(ctx context.Context, runnable utils.Runnable, username string) (*response.User, error)

	CheckPhoneAlreadyExists(ctx context.Context, value string) (bool, error)
	CheckEmailAlreadyExists(ctx context.Context, runnable utils.Runnable, value string) (*response.User, error)
	CheckUsernameAlreadyExists(ctx context.Context, runnable utils.Runnable, username string) (*response.User, error)

	CreateGoatInviteCode(ctx context.Context, code string) (string, error)
	GetInviteCodeStatus(ctx context.Context, code string) (*int, bool, error)
	GetNextAvailableUsername(ctx context.Context, firstName string, lastName string) (*string, error)
	UseInviteCode(ctx context.Context, tx pgx.Tx, code string, userID int) error
	UpsertStripeAccountIDByEmail(ctx context.Context, email string, stripeAccountID string) error
	UpsertUser(ctx context.Context, tx pgx.Tx, user *response.User) (*response.User, error)
	UpdateCurrentUser(ctx context.Context, req request.UpdateUser) (*response.User, error)
	UpsertGoatChatsPrice(ctx context.Context, tx pgx.Tx, priceInSmallestDenom int64, currency string, goatUserID int) error
	UploadProfileAvatar(ctx context.Context, req request.UploadProfileAvatar) (*response.UploadProfileAvatar, error)
	UpsertBio(ctx context.Context, req request.UpsertBio) (*response.UpsertBio, error)
	UpdatePhone(ctx context.Context, req request.UpdatePhone) (*response.UpdatePhone, error)
	GetGoatUsers(ctx context.Context, req request.GetGoatUsers) (*response.GetGoatUsers, error)
	GetUserByStripeAccountID(ctx context.Context, runnable utils.Runnable, stripeAccountID string) (*response.User, error)
	GetUserByStripeID(ctx context.Context, runnable utils.Runnable, stripeID string) (*response.User, error)
	UpgradeUserToGoat(ctx context.Context, formattedNumber string, req request.SignupSMSGoat) (*response.AuthSuccess, error)
	IsContact(ctx context.Context, contactID int) (*response.IsContact, error)
	GetInviteCodeStatuses(ctx context.Context) (*response.GetInviteCodeStatuses, error)

	BlockUser(ctx context.Context, blockUserID int) error
	UnblockUser(ctx context.Context, unblockUserID int) error
	IsBlocked(ctx context.Context, userID int) (bool, error)
}
