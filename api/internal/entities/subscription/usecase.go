package subscription

import (
	"context"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
)

type Usecase interface {
	SubscribeCurrUserToGoat(c context.Context, req request.SubscribeCurrUserToGoat) (*response.UserSubscription, error)
	GetUserSubscriptions(ctx context.Context, req request.GetUserSubscriptions) (*response.GetUserSubscriptions, error)
	GetUserSubscriptionByGoatID(ctx context.Context, userID int, goatUserID int) (*response.UserSubscription, error)
	UnsubscribeCurrUserFromGoat(ctx context.Context, req request.UnsubscribeCurrUserFromGoat) (*response.UserSubscription, error)

	UpsertUserSubscription(ctx context.Context, runnable utils.Runnable, subscription response.UserSubscription) error
	GetSubscriptionPrices(ctx context.Context, req request.GetSubscriptionPrices) (*response.GoatSubscriptionInfo, error)

	UpsertGoatSubscriptionInfo(ctx context.Context, priceInSmallestDenom int64, currency string) error
	UpsertGoatSubscriptionInfoTx(ctx context.Context, tx pgx.Tx, goatUser *response.User, priceInSmallestDenom int64, currency string) error

	CheckUserSubscribedToGoat(ctx context.Context, req request.CheckUserSubscribedToGoat) (*response.CheckUserSubscribedToGoat, error)

	SubscribePaidGroup(ctx context.Context, tx pgx.Tx, channelID string, user *response.User) (*response.PaidGroupChatSubscription, error)
	UnsubscribePaidGroup(ctx context.Context, tx pgx.Tx, channelID string, user *response.User) (*response.PaidGroupChatSubscription, error)
	UpsertPaidGroupSubscription(ctx context.Context, runnable utils.Runnable, subscription *response.PaidGroupChatSubscription) error
	GetMyPaidGroupSubscriptions(ctx context.Context, req request.GetMyPaidGroupSubscriptions) (*response.GetMyPaidGroupSubscriptions, error)
	GetMyPaidGroupSubscription(ctx context.Context, req request.GetMyPaidGroupSubscription) (*response.GetMyPaidGroupSubscription, error)

	BatchUnsubscribePaidGroup(ctx context.Context, req cloudtasks.StripePaidGroupUnsubscribeTask) error
	DeletePaidGroupProduct(ctx context.Context, runnable utils.Runnable, channelID string) error
	DeletePaidGroupSubscription(ctx context.Context, runnable utils.Runnable, channelID string, userID int) error
}
