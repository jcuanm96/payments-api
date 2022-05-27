package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) UpsertGoatSubscriptionInfo(ctx context.Context, tx pgx.Tx, goatUserID int, tierName string, priceInSmallestDenom int64, currency string, stripeProductID string) error {
	query, args, squirrelErr := squirrel.Insert("subscription.tiers").
		Columns(
			"goat_user_id",
			"price_in_smallest_denom",
			"currency",
			"tier_name",
			"stripe_product_id",
		).
		Values(
			goatUserID,
			priceInSmallestDenom,
			currency,
			tierName,
			stripeProductID,
		).
		Suffix(`
			ON CONFLICT ON CONSTRAINT tiers_goat_user_id_tier_name_key
			DO UPDATE SET 
			price_in_smallest_denom = ?,
			currency = ?,
			stripe_product_id = ?`,
			priceInSmallestDenom,
			currency,
			stripeProductID,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := tx.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
