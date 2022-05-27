package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) ListUnpaidProviders(ctx context.Context) ([]response.ProviderPaymentInfo, error) {
	query, args, squirrelErr := squirrel.Select(
		"users.id",
		"users.first_name",
		"users.last_name",

		"COALESCE(balances.available_balance, 0)",
		"COALESCE(balances.currency, 'usd')",
		"COALESCE(balances.last_payout_ts, to_timestamp('2000-01-01', 'YYYY-MM-DD') at time zone 'Etc/UTC')", // Arbitrary low date
		"COALESCE(balances.last_paid_payout_period_id, -1)",

		"COALESCE(bank_info.bank_name, '')",
		"COALESCE(bank_info.account_number, '')",
		"COALESCE(bank_info.routing_number, '')",
		"COALESCE(bank_info.account_type, '')",
		"COALESCE(bank_info.account_holder_name, '')",
		"COALESCE(bank_info.account_holder_type, '')",
		"COALESCE(bank_info.currency, '')",
		"COALESCE(bank_info.country, '')",

		"COALESCE(addresses.street_1, '')",
		"COALESCE(addresses.street_2, '')",
		"COALESCE(addresses.city, '')",
		"COALESCE(addresses.state, '')",
		"COALESCE(addresses.postal_code, '')",
		"COALESCE(addresses.country, '')",
	).
		From("wallet.balances balances").
		RightJoin("core.users users ON users.id = balances.provider_user_id").
		LeftJoin("wallet.bank_info bank_info ON bank_info.user_id = balances.provider_user_id").
		LeftJoin("wallet.billing_addresses addresses ON bank_info.user_id = addresses.user_id").
		Where("users.user_type = 'GOAT'").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}

	providers := []response.ProviderPaymentInfo{}

	defer rows.Close()
	for rows.Next() {
		provider := response.ProviderPaymentInfo{}
		bank := response.Bank{}
		scanErr := rows.Scan(
			&provider.ID,
			&provider.FirstName,
			&provider.LastName,

			&provider.BalanceOwed,
			&provider.Currency,
			&provider.LastPayoutTS,
			&provider.LastPayoutPeriodID,

			&bank.BankName,
			&bank.AccountNumber,
			&bank.RoutingNumber,
			&bank.AccountType,
			&bank.AccountHolderName,
			&bank.AccountHolderType,
			&bank.Currency,
			&bank.Country,

			&bank.BillingAddress.Street1,
			&bank.BillingAddress.Street2,
			&bank.BillingAddress.City,
			&bank.BillingAddress.State,
			&bank.BillingAddress.PostalCode,
			&bank.BillingAddress.Country,
		)
		if scanErr != nil {
			return nil, scanErr
		}

		provider.Bank = bank
		providers = append(providers, provider)
	}

	return providers, nil
}
