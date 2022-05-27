package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) CheckUserByFieldAlreadyExists(ctx context.Context, fieldName, fieldValue string) (bool, error) {
	res := 0
	query, args, squirrelErr := squirrel.Select(
		"id",
	).From("core.users").Where(fieldName+" = ?", fieldValue).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return false, squirrelErr
	}
	row := s.MasterNode().QueryRow(ctx, query, args...)
	scanErr := row.Scan(
		&res,
	)
	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return res > 0, nil
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return res > 0, scanErr
	}
	return res > 0, nil
}
