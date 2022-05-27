package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) ListPayoutHistory(ctx context.Context, req request.ListPayoutHistory, limit int64) ([]response.PayoutHistoryDatum, error) {
	currency := "usd"
	args := []interface{}{req.GoatUserID, currency, req.CursorPayPeriodID, limit}
	query := `
		WITH balance_owed_per_period AS (
			SELECT
				pay_period_id,
				SUM(amount - (stripe_fee + vama_fee)) AS total
			FROM wallet.ledger 
			WHERE 
				source_type = 'charge' AND
				provider_user_id = $1 AND
				currency = $2
			GROUP BY pay_period_id
		), balance_paid_per_period AS (
			SELECT
				pay_period_id,
				SUM(amount) AS total
			FROM wallet.ledger 
			WHERE 
				source_type = 'PAYOUT' AND
				provider_user_id = $1 AND
				currency = $2
			GROUP BY pay_period_id
		)

		SELECT 
			payout_periods.id,
			payout_periods.start_ts,
			payout_periods.end_ts,
			COALESCE(balance_owed.total,0) AS total_balance_owed,
			COALESCE(balance_paid.total, 0) AS total_balance_paid
		FROM wallet.payout_periods payout_periods
		LEFT JOIN balance_owed_per_period balance_owed ON payout_periods.id = balance_owed.pay_period_id
		LEFT JOIN balance_paid_per_period balance_paid ON payout_periods.id = balance_paid.pay_period_id
		WHERE payout_periods.id < $3
		ORDER BY payout_periods.id DESC
		LIMIT $4
	`

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}

	payoutHistory := []response.PayoutHistoryDatum{}

	defer rows.Close()

	for rows.Next() {
		currPayoutHistoryDatum := response.PayoutHistoryDatum{}
		currPayoutPeriod := response.PayoutPeriod{}
		scanErr := rows.Scan(
			&currPayoutPeriod.ID,
			&currPayoutPeriod.StartTS,
			&currPayoutPeriod.EndTS,
			&currPayoutHistoryDatum.BalanceOwed,
			&currPayoutHistoryDatum.BalancePaid,
		)

		if scanErr != nil {
			vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
			return nil, scanErr
		}

		currPayoutHistoryDatum.PayoutPeriod = currPayoutPeriod

		payoutHistory = append(payoutHistory, currPayoutHistoryDatum)
	}

	return payoutHistory, nil
}
