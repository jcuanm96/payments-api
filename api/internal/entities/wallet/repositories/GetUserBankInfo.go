package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetUserBankInfo(ctx context.Context, userID int) ([]response.PaymentMethod, error) {
	queryString, args, queryErr := squirrel.Select(
		"bank_info.bank_name",
		"RIGHT(bank_info.account_number, 4)", // get last 4 digits
		"bank_info.account_type",
		"bank_info.account_holder_name",
		"bank_info.account_holder_type",
		"bank_info.currency",
		"bank_info.country",

		"COALESCE(addresses.street_1, '')",
		"COALESCE(addresses.street_2, '')",
		"COALESCE(addresses.city, '')",
		"COALESCE(addresses.state, '')",
		"COALESCE(addresses.postal_code, '')",
		"COALESCE(addresses.country, '')",
	).
		From("wallet.bank_info bank_info").
		LeftJoin("wallet.billing_addresses addresses ON bank_info.user_id = addresses.user_id").
		Where(`bank_info.user_id = ?`, userID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if queryErr != nil {
		return nil, queryErr
	}

	rows, sqlErr := s.MasterNode().Query(ctx, queryString, args...)
	if sqlErr != nil {
		return nil, sqlErr
	}

	defer rows.Close()
	banks := []response.PaymentMethod{}

	for rows.Next() {
		currBank := response.Bank{}
		scanErr := rows.Scan(
			&currBank.BankName,
			&currBank.AccountNumber,
			&currBank.AccountType,
			&currBank.AccountHolderName,
			&currBank.AccountHolderType,
			&currBank.Currency,
			&currBank.Country,

			&currBank.BillingAddress.Street1,
			&currBank.BillingAddress.Street2,
			&currBank.BillingAddress.City,
			&currBank.BillingAddress.State,
			&currBank.BillingAddress.PostalCode,
			&currBank.BillingAddress.Country,
		)

		if scanErr != nil {
			vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, queryString, args))
			return nil, scanErr
		}

		banks = append(banks, response.PaymentMethod{Bank: &currBank})
	}

	return banks, nil
}
