package repositories

import (
	"context"
	"errors"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/jackc/pgx/v4"
)

func (s *repository) UpdateBalancePayout(ctx context.Context, tx pgx.Tx, providerID int, amountPaidNegative int64, payPeriodID int) error {
	if amountPaidNegative >= 0 {
		return errors.New("amountPaidNegative was not negative")
	}
	currency := constants.DEFAULT_CURRENCY
	query := `
		UPDATE wallet.balances
		SET
			available_balance = available_balance + $1,
			last_payout_ts = now(),
			last_paid_payout_period_id = $4
		WHERE
			provider_user_id = $2 AND
			currency = $3
	`

  args := []interface{}{amountPaidNegative, providerID, currency, payPeriodID}
	result, queryErr := tx.Exec(ctx, query, args...)
	if queryErr != nil {
		return queryErr
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected < 1 {
		return constants.ErrNoRowsAffected
	}

	return nil
}
