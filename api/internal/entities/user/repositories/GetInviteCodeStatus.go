package repositories

import (
	"context"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetInviteCodeStatus(ctx context.Context, code string) (*int, bool, error) {
	var res *int
	query, args, squirrelErr := squirrel.Select(
		"used_by",
	).From("core.goat_invite_codes").Where("invite_code = ?", code).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return nil, false, squirrelErr
	}
	row := s.MasterNode().QueryRow(ctx, query, args...)
	scanErr := row.Scan(
		&res,
	)
	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return res, false, nil
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return res, false, scanErr
	}
	return res, true, nil
}
