package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func EnableGoatPostNotifications(ctx context.Context, runnable utils.Runnable, userID, goatID int) error {
	query, args, squirrelErr := squirrel.Insert("feed.post_notifications").
		Columns(
			"user_id",
			"goat_user_id",
		).
		Values(
			userID,
			goatID,
		).
		Suffix(`ON conflict (user_id, goat_user_id) do nothing`).
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

func (s *repository) EnableGoatPostNotifications(ctx context.Context, runnable utils.Runnable, userID int, goatID int) error {
	return EnableGoatPostNotifications(ctx, runnable, userID, goatID)
}
