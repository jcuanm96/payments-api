package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) DeleteComment(ctx context.Context, commentID int, userID int) error {
	sql, args, sqlErr := squirrel.Update("feed.post_comments").
		Set("deleted_at", "NOW()").
		Where("id = ?", commentID).
		Where("user_id = ?", userID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if sqlErr != nil {
		vlog.Errorf(ctx, utils.SqlErrLogMsg(sqlErr, sql, args))
		return sqlErr
	}
	result, execErr := s.MasterNode().Exec(ctx, sql, args...)
	if execErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(execErr, sql, args))
		return execErr
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected < 1 {
		return constants.ErrNoRowsAffected
	}

	return nil
}
