package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (s *repository) GetPendingTransactionHistory(ctx context.Context, userID int, lastTransactionID int64, limit uint64) ([]response.TransactionItem, error) {
	query := squirrel.Select(
		"pending.id",
		"pending.amount",
		"pending.currency",
		"pending.stripe_created_ts",

		"provider.id",
		"provider.first_name",
		"provider.last_name",
		"provider.phone_number",
		"provider.country_code",
		"provider.email",
		"provider.username",
		"provider.user_type",
		"provider.profile_avatar",
	).
		From("wallet.pending_transactions pending").
		Join("core.users provider ON provider.id = pending.provider_user_id").
		Where("pending.customer_user_id = ?", userID)

	if lastTransactionID > 0 {
		query = query.Where("pending.id < ?", lastTransactionID)
	}

	queryStr, args, queryErr := query.
		OrderBy("pending.id DESC").
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
		scanErr := rows.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.CreatedAt,

			&provider.ID,
			&provider.FirstName,
			&provider.LastName,
			&provider.Phonenumber,
			&provider.CountryCode,
			&provider.Email,
			&provider.Username,
			&provider.Type,
			&provider.ProfileAvatar,
		)
		if scanErr != nil {
			return nil, scanErr
		}

		transaction.User = &provider
		transaction.Type = constants.TXN_HISTORY_PENDING

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
