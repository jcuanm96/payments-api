package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) UpsertBank(ctx context.Context, tx pgx.Tx, userID int, req request.UpsertBank) error {
	query, args, squirrelErr := squirrel.Insert("wallet.bank_info").
		Columns(
			"user_id",
			"bank_name",
			"account_number",
			"routing_number",
			"account_type",
			"account_holder_name",
			"account_holder_type",
			"currency",
			"country",
		).
		Values(
			userID,
			req.BankName,
			req.AccountNumber,
			req.RoutingNumber,
			req.AccountType,
			req.AccountHolderName,
			req.AccountHolderType,
			req.Currency,
			req.Country,
		).
		Suffix(`
			ON conflict (user_id)
			DO UPDATE SET
			bank_name = ?,
			account_number = ?,
			routing_number = ?,
			account_type = ?,
			account_holder_name = ?,
			account_holder_type = ?,
			currency = ?,
			country = ?`,
			req.BankName,
			req.AccountNumber,
			req.RoutingNumber,
			req.AccountType,
			req.AccountHolderName,
			req.AccountHolderType,
			req.Currency,
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
