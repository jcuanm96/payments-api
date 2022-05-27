package subscription

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
)

type Repository interface {
	baserepo.BaseRepository
	GetUserSubscriptionByGoatID(ctx context.Context, userID int, goatUserID int) (*response.UserSubscription, error)
	GetUserSubscriptions(ctx context.Context, userID int, cursorID int, limit uint64) ([]response.UserSubscription, *int, error)
	UpsertUserSubscription(ctx context.Context, runnable utils.Runnable, subscription response.UserSubscription) error

	GetGoatSubscriptionInfo(ctx context.Context, goatUserID int) (*response.GoatSubscriptionInfo, error)
	UpsertGoatSubscriptionInfo(ctx context.Context, tx pgx.Tx, goatUserID int, tierName string, priceInSmallestDenom int64, currency string, stripeProductID string) error

	GetPaidGroupProductInfo(ctx context.Context, runnable utils.Runnable, channelID string) (*PaidGroupChatInfo, error)

	GetPaidGroupSubscription(ctx context.Context, runnable utils.Runnable, userID int, channelID string) (*response.PaidGroupChatSubscription, error)
	UpsertPaidGroupSubscription(ctx context.Context, runnable utils.Runnable, subscription *response.PaidGroupChatSubscription) error
	GetMyPaidGroupSubscriptions(ctx context.Context, userID int, cursorID int, limit uint64) ([]response.PaidGroupChatSubscription, error)
	DeletePaidGroupProduct(ctx context.Context, runnable utils.Runnable, channelID string) error
	DeletePaidGroupSubscription(ctx context.Context, runnable utils.Runnable, channelID string, userID int) error
}
