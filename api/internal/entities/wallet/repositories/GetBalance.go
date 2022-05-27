package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetBalance(ctx context.Context, providerUserID int, currency string) (*int64, error) {
	query, args, squirrelErr := squirrel.Select(
		`available_balance`,
	).
		From("wallet.balances").
		Where("provider_user_id = ?", providerUserID).
		Where("currency = ?", currency).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if squirrelErr != nil {
		return nil, squirrelErr
	}
	row := s.MasterNode().QueryRow(ctx, query, args...)

	var availableBalance int64
	scanErr := row.Scan(
		&availableBalance,
	)
	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			availableBalance = 0
			return &availableBalance, nil
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return nil, scanErr
	}

	return &availableBalance, nil

}
