package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetUserSubscriptionByGoatID(ctx context.Context, userID int, goatUserID int) (*response.UserSubscription, error) {
	subscription := response.UserSubscription{}
	query, args, squirrelErr := squirrel.Select(
		"id",
		"current_period_end",
		"user_id",
		"stripe_subscription_id",
		"goat_user_id",
		"tier_id",
		"is_renewing",
	).From("subscription.user_subscriptions").
		Where("user_id = ?", userID).
		Where("goat_user_id = ?", goatUserID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	row := s.MasterNode().QueryRow(ctx, query, args...)

	scanErr := row.Scan(
		&subscription.ID,
		&subscription.CurrentPeriodEnd,
		&subscription.UserID,
		&subscription.StripeSubscriptionID,
		&subscription.GoatUser.ID,
		&subscription.TierID,
		&subscription.IsRenewing,
	)

	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return nil, nil
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return nil, scanErr
	}

	return &subscription, scanErr
}
