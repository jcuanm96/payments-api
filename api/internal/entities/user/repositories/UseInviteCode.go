package repositories

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) UseInviteCode(ctx context.Context, tx pgx.Tx, code string, userID int) error {
	query, args, squirrelErr := squirrel.Update("core.goat_invite_codes").
		Set("used_by", userID).
		Set("updated_at", time.Now()).
		Where("invite_code = ?", code).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}
	_, queryErr := tx.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
