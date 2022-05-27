package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func UpsertPaidGroupSubscription(ctx context.Context, runnable utils.Runnable, subscription *response.PaidGroupChatSubscription) error {
	query, args, squirrelErr := squirrel.Insert("subscription.paid_group_chat_subscriptions").
		Columns(
			"user_id",
			"goat_user_id",
			"stripe_subscription_id",
			"sendbird_channel_id",
			"current_period_end",
			"is_renewing",
		).
		Values(
			subscription.UserID,
			subscription.GoatUser.ID,
			subscription.StripeSubscriptionID,
			subscription.ChannelID,
			subscription.CurrentPeriodEnd,
			subscription.IsRenewing,
		).
		Suffix(`
		ON CONFLICT (user_id, sendbird_channel_id)
		DO UPDATE SET
			stripe_subscription_id = ?,
			goat_user_id = ?,
			current_period_end = ?,
			is_renewing = ?`,
			subscription.StripeSubscriptionID,
			subscription.GoatUser.ID,
			subscription.CurrentPeriodEnd,
			subscription.IsRenewing,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := runnable.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}
	return nil
}

func (s *repository) UpsertPaidGroupSubscription(ctx context.Context, runnable utils.Runnable, subscription *response.PaidGroupChatSubscription) error {
	return UpsertPaidGroupSubscription(ctx, runnable, subscription)
}
