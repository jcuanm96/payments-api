package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateGoatInviteCodeDB(ctx context.Context, code string, db *pgxpool.Pool) error {
	query, args, squirrelErr := squirrel.Insert("core.goat_invite_codes").
		Columns(
			"invite_code",
		).
		Values(
			code,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}
	_, queryErr := db.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}

func (s *repository) CreateGoatInviteCode(ctx context.Context, code string) error {
	return CreateGoatInviteCodeDB(ctx, code, s.MasterNode())
}
