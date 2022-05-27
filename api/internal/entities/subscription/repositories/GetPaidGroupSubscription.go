package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetPaidGroupSubscription(ctx context.Context, runnable utils.Runnable, userID int, channelID string) (*response.PaidGroupChatSubscription, error) {
	return GetPaidGroupSubscription(ctx, runnable, userID, channelID)
}

func GetPaidGroupSubscription(ctx context.Context, runnable utils.Runnable, userID int, channelID string) (*response.PaidGroupChatSubscription, error) {
	query, args, squirrelErr := squirrel.Select(
		"id",
		"current_period_end",
		"user_id",
		"stripe_subscription_id",
		"goat_user_id",
		"sendbird_channel_id",
		"is_renewing",
	).From("subscription.paid_group_chat_subscriptions").
		Where("user_id = ?", userID).
		Where("sendbird_channel_id = ?", channelID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	row := runnable.QueryRow(ctx, query, args...)

	subscription := response.PaidGroupChatSubscription{}
	scanErr := row.Scan(
		&subscription.ID,
		&subscription.CurrentPeriodEnd,
		&subscription.UserID,
		&subscription.StripeSubscriptionID,
		&subscription.GoatUser.ID,
		&subscription.ChannelID,
		&subscription.IsRenewing,
	)

	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return nil, nil
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return nil, scanErr
	}

	return &subscription, nil
}
