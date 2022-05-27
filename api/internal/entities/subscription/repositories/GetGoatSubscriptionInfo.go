package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetGoatSubscriptionInfo(ctx context.Context, goatUserID int) (*response.GoatSubscriptionInfo, error) {
	tier := response.GoatSubscriptionInfo{}
	query, args, squirrelErr := squirrel.Select(
		"id",
		"goat_user_id",
		"price_in_smallest_denom",
		"currency",
		"tier_name",
		"stripe_product_id",
	).From("subscription.tiers").
		Where("tier_name = ?", constants.DEFAULT_TIER_NAME).
		Where("goat_user_id = ?", goatUserID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	row := s.MasterNode().QueryRow(ctx, query, args...)
	scanErr := row.Scan(
		&tier.ID,
		&tier.GoatUserID,
		&tier.PriceInSmallestDenom,
		&tier.Currency,
		&tier.TierName,
		&tier.StripeProductID,
	)

	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return nil, nil
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return nil, scanErr
	}

	return &tier, nil
}
