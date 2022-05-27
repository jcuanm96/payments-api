package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) UpdateSubscriptionPrice(ctx context.Context, goatID int, tierName string, price int, currency string) error {
	query, args, squirrelErr := squirrel.Update("subscription.tiers").
		Set("price_in_smallest_denom", price).
		Set("currency", currency).
		Where("goat_user_id = ?", goatID).
		Where("tier_name = ?", tierName).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := s.MasterNode().Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
