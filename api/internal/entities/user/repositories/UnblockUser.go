package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) UnblockUser(ctx context.Context, tx pgx.Tx, currUserID int, unblockUserID int) error {
	query, args, squirrelErr := squirrel.Delete("core.user_blocks").
		Where("user_id = ?", currUserID).
		Where("blocked_user_id = ?", unblockUserID).
		PlaceholderFormat(squirrel.Dollar).ToSql()

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
