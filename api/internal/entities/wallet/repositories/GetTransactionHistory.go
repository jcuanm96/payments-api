package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/stripe/stripe-go/v72"
)

func (s *repository) GetTransactionHistory(ctx context.Context, userID int, lastTransactionID int64, limit uint64) ([]response.TransactionItem, error) {
	query := squirrel.Select(
		"ledger.id",
		"ledger.amount",
		"ledger.currency",
		"ledger.stripe_created_ts",
		"ledger.source_type",
		"COALESCE(ledger.vama_fee, 0)",
		"COALESCE(ledger.stripe_fee, 0)",
		"CAST(extract(epoch FROM ledger.created_at) AS BIGINT) AS created_at",

		"provider.id",
		"provider.first_name",
		"provider.last_name",
		"provider.phone_number",
		"provider.country_code",
		"provider.email",
		"provider.username",
		"provider.user_type",
		"provider.profile_avatar",

		"COALESCE(customer.id, 0)",
		"COALESCE(customer.first_name, '')",
		"COALESCE(customer.last_name, '')",
		"COALESCE(customer.phone_number, '')",
		"COALESCE(customer.country_code, '')",
		"COALESCE(customer.email, '')",
		"COALESCE(customer.username, '')",
		"COALESCE(customer.user_type, '')",
		"COALESCE(customer.profile_avatar, '')",
	).
		From("wallet.ledger ledger").
		Join("core.users provider ON provider.id = ledger.provider_user_id").
		LeftJoin("core.users customer ON customer.id = ledger.customer_user_id").
		Where("(ledger.provider_user_id = ? OR ledger.customer_user_id = ?)", userID, userID)

	if lastTransactionID > 0 {
		query = query.Where("ledger.id < ?", lastTransactionID)
	}

	queryStr, args, queryErr := query.
		OrderBy("ledger.id DESC").
		Limit(limit).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if queryErr != nil {
		return nil, queryErr
	}

	rows, sqlErr := s.MasterNode().Query(ctx, queryStr, args...)
	if sqlErr != nil {
		return nil, sqlErr
	}

	defer rows.Close()
	transactions := []response.TransactionItem{}

	for rows.Next() {
		transaction := response.TransactionItem{}

		provider := response.User{}
		customer := response.User{}
		var sourceType string
		var vamaFee int64
		var stripeFee int64
		var vamaCreatedAt int64

		scanErr := rows.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.CreatedAt,
			&sourceType,
			&vamaFee,
			&stripeFee,
			&vamaCreatedAt,

			&provider.ID,
			&provider.FirstName,
			&provider.LastName,
			&provider.Phonenumber,
			&provider.CountryCode,
			&provider.Email,
			&provider.Username,
			&provider.Type,
			&provider.ProfileAvatar,

			&customer.ID,
			&customer.FirstName,
			&customer.LastName,
			&customer.Phonenumber,
			&customer.CountryCode,
			&customer.Email,
			&customer.Username,
			&customer.Type,
			&customer.ProfileAvatar,
		)
		if scanErr != nil {
			return nil, scanErr
		}

		if sourceType == constants.TXN_HISTORY_PAYOUT {
			transaction.User = &provider
			transaction.Type = constants.TXN_HISTORY_PAYOUT
			transaction.CreatedAt = vamaCreatedAt
			transaction.Amount = transaction.Amount * -1
		} else if sourceType == string(stripe.BalanceTransactionSourceTypeCharge) {
			transaction.Fees = vamaFee + stripeFee
			if provider.ID == userID {
				transaction.User = &customer
				transaction.Type = constants.TXN_HISTORY_INCOMING
			} else if customer.ID == userID {
				transaction.User = &provider
				transaction.Type = constants.TXN_HISTORY_OUTGOING
			}

			if customer.ID == constants.VAMA_USER_ID {
				transaction.CreatedAt = vamaCreatedAt
			}
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
