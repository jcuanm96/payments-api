package user

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
)

type Repository interface {
	baserepo.BaseRepository

	CheckUserByFieldAlreadyExists(ctx context.Context, fieldName, fieldValue string) (bool, error)
	IsLinkSuffixTaken(ctx context.Context, runnable utils.Runnable, suffix string) (*LinkSuffixTaken, error)

	GetUserByEmail(ctx context.Context, runnable utils.Runnable, email string) (*response.User, error)
	GetUserByUsername(ctx context.Context, runnable utils.Runnable, username string) (*response.User, error)
	GetUserByPhone(ctx context.Context, runnable utils.Runnable, phone string) (*response.User, error)
	GetUserByID(ctx context.Context, id int) (*response.User, error)
	GetUserByUUID(ctx context.Context, runnable utils.Runnable, uuid string) (*response.User, error)
	GetUser(ctx context.Context, runnable utils.Runnable, prefict string, params ...interface{}) (*response.User, error)
	CheckUserDeleted(ctx context.Context, userUUID string) (bool, error)
	UpdateUser(ctx context.Context, runnable utils.Runnable, userID int, item request.UpdateUser) error
	SoftDeleteUser(ctx context.Context, userID int) error
	HardDeleteUser(ctx context.Context, userID int) error
	UndeleteUser(ctx context.Context, userID int) error
	UpsertUser(ctx context.Context, tx pgx.Tx, user *response.User) (*response.User, error)
	UpdateUserProfileAvatar(ctx context.Context, runnable utils.Runnable, userID int, profileAvatarFilename string) error
	GetBioByID(ctx context.Context, userID int) (string, error)
	UpsertBio(ctx context.Context, userID int, bio string) error
	UpsertStripeAccountIDByEmail(ctx context.Context, email string, profileAvatarFilename string) error
	UpsertGoatChatsPrice(ctx context.Context, tx pgx.Tx, priceInSmallestDenom int64, currency string, goatUserID int) error
	ReportUser(ctx context.Context, reporterUserID int, reportedUserID int, description string) error

	CreateGoatInviteCode(ctx context.Context, code string) error
	GetMyInvites(ctx context.Context, userID int) ([]response.Invite, error)
	GetInviteCodeStatus(ctx context.Context, code string) (*int, bool, error)
	GetNextAvailableUsername(ctx context.Context, firstName string, lastName string) (*string, error)
	UseInviteCode(ctx context.Context, tx pgx.Tx, code string, userID int) error
	UpdatePhone(ctx context.Context, runnable utils.Runnable, userID int, phone string) error
	GetGoatUsers(ctx context.Context, currUserID int, req request.GetGoatUsers) ([]response.GetGoatUsersListItem, error)
	GetGoatProfile(ctx context.Context, goatID, currUserID int) (*response.GetGoatProfile, error)
	GetUserByStripeAccountID(ctx context.Context, runnable utils.Runnable, stripeAccountID string) (*response.User, error)
	GetUserByStripeID(ctx context.Context, runnable utils.Runnable, stripeID string) (*response.User, error)

	IsContact(ctx context.Context, userID, contactID int) (bool, error)

	BlockUser(ctx context.Context, runnable utils.Runnable, currUserID int, blockUserID int) error
	UnblockUser(ctx context.Context, tx pgx.Tx, currUserID int, unblockUserID int) error
	IsBlocked(ctx context.Context, currUserID int, userID int) (bool, error)

	CensorUserData(ctx context.Context, exec utils.Executable, userID int) error
	ClearFollows(ctx context.Context, exec utils.Executable, userID int) error
	ClearFCMTokens(ctx context.Context, exec utils.Executable, userID int) error
	ClearAuthTokens(ctx context.Context, exec utils.Executable, userID int) error
}
