package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) DisableGoatPostNotifications(ctx context.Context, runnable utils.Runnable, userID, goatID int) error {
	query, args, squirrelErr := squirrel.Delete("feed.post_notifications").
		Where("user_id = ?", userID).
		Where("goat_user_id = ?", goatID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := runnable.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
