package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (s *repository) GetPayPeriodPerTS(ctx context.Context, ts int64) (*response.PayoutPeriod, error) {
	queryString, args, queryErr := squirrel.Select(
		"id",
		"start_ts",
		"end_ts",
	).
		From("wallet.payout_periods").
		Where("start_ts <= ?", ts).
		Where("end_ts > ?", ts).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if queryErr != nil {
		return nil, queryErr
	}

	row, sqlErr := s.MasterNode().Query(ctx, queryString, args...)
	if sqlErr != nil {
		return nil, sqlErr
	}

	defer row.Close()
	hasNext := row.Next()

	if !hasNext {
		return nil, constants.ErrNotFound
	}

	payPeriod := response.PayoutPeriod{}
	scanErr := row.Scan(
		&payPeriod.ID,
		&payPeriod.StartTS,
		&payPeriod.EndTS,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	return &payPeriod, nil

}
