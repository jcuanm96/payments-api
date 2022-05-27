package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) UpsertGoatChatsPrice(ctx context.Context, tx pgx.Tx, priceInSmallestDenom int64, currency string, goatUserID int) error {
	query, args, squirrelErr := squirrel.Insert("product.goat_chats").
		Columns(
			"price_in_smallest_denom",
			"currency",
			"goat_user_id",
		).
		Values(
			priceInSmallestDenom,
			currency,
			goatUserID,
		).
		Suffix(`ON conflict (goat_user_id)
			do update SET 
			price_in_smallest_denom = ?,
			currency = ?`,
			priceInSmallestDenom,
			currency,
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
