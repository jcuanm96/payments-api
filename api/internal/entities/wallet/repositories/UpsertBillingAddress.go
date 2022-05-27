package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) UpsertBillingAddress(ctx context.Context, tx pgx.Tx, userID int, req request.BillingAddress) error {
	query, args, squirrelErr := squirrel.Insert("wallet.billing_addresses").
		Columns(
			"user_id",
			"street_1",
			"street_2",
			"city",
			"state",
			"postal_code",
			"country",
		).
		Values(
			userID,
			req.Street1,
			req.Street2,
			req.City,
			req.State,
			req.PostalCode,
			req.Country,
		).
		Suffix(`
			ON conflict (user_id)
			DO UPDATE SET
			street_1 = ?,
			street_2 = ?,
			city = ?,
			state = ?,
			postal_code = ?,
			country = ?`,
			req.Street1,
			req.Street2,
			req.City,
			req.State,
			req.PostalCode,
			req.Country,
		).
		PlaceholderFormat(squirrel.Dollar).ToSql()
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
