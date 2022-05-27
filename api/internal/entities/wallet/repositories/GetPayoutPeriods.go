package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetPayoutPeriods(ctx context.Context, req request.GetPayoutPeriods) ([]response.PayoutPeriod, error) {
	query, args, squirrelErr := squirrel.Select(
		"id",
		"start_ts",
		"end_ts",
	).
		From("wallet.payout_periods").
		Where("id > ?", req.CursorID).
		OrderBy("id ASC").
		Limit(uint64(req.Limit)).
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

	defer rows.Close()
	payoutPeriods := []response.PayoutPeriod{}
	for rows.Next() {
		payoutPeriod := response.PayoutPeriod{}
		scanErr := rows.Scan(
			&payoutPeriod.ID,
			&payoutPeriod.StartTS,
			&payoutPeriod.EndTS,
		)

		if scanErr != nil {
			return nil, scanErr
		}

		payoutPeriods = append(payoutPeriods, payoutPeriod)
	}
	return payoutPeriods, nil
}
