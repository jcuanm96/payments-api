package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func BlockUser(ctx context.Context, runnable utils.Runnable, currUserID int, blockUserID int) error {
	query, args, squirrelErr := squirrel.Insert("core.user_blocks").
		Columns(
			"user_id",
			"blocked_user_id",
		).
		Values(
			currUserID,
			blockUserID,
		).
		Suffix(`ON conflict (user_id, blocked_user_id) do nothing`).
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

func (s *repository) BlockUser(ctx context.Context, runnable utils.Runnable, userID int, blockUserID int) error {
	return BlockUser(ctx, runnable, userID, blockUserID)
}
