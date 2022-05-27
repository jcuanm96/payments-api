package repositories

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) UpsertUserSubscription(ctx context.Context, runnable utils.Runnable, subscription response.UserSubscription) error {
	if subscription.GoatUser.ID == 0 {
		return errors.New("GoatUser.ID passed as 0 to UpsertNewSubscription")
	}
	query, args, squirrelErr := squirrel.Insert("subscription.user_subscriptions").
		Columns(
			"user_id",
			"goat_user_id",
			"stripe_subscription_id",
			"tier_id",
			"current_period_end",
			"is_renewing",
		).
		Values(
			subscription.UserID,
			subscription.GoatUser.ID,
			subscription.StripeSubscriptionID,
			subscription.TierID,
			subscription.CurrentPeriodEnd,
			subscription.IsRenewing,
		).
		Suffix(`
			ON CONFLICT (user_id, goat_user_id)
			DO UPDATE SET
				stripe_subscription_id = ?,
				tier_id = ?,
				current_period_end = ?,
				is_renewing = ?`,
			subscription.StripeSubscriptionID,
			subscription.TierID,
			subscription.CurrentPeriodEnd,
			subscription.IsRenewing,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	result, execContextErr := runnable.Exec(ctx, query, args...)
	if execContextErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(execContextErr, query, args))
		return execContextErr
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected < 1 {
		return constants.ErrNoRowsAffected
	}

	return nil
}
