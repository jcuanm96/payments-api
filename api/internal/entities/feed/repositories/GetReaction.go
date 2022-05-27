package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetReaction(ctx context.Context, userID int, postID int, runnable utils.Runnable) (*string, error) {
	query, _, selectErr := squirrel.Select(
		"type",
	).
		From("feed.post_reactions").
		Where("post_id = ?").
		Where("user_id = ?").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if selectErr != nil {
		return nil, selectErr
	}

	args := []interface{}{postID, userID}
	row, queryErr := runnable.Query(ctx, query, args...)
	if queryErr != nil {
		return nil, queryErr
	}

	defer row.Close()
	hasNext := row.Next()

	if !hasNext {
		return &appconfig.Config.Vote.Nil, nil
	}

	var reaction string
	scanErr := row.Scan(
		&reaction,
	)

	if scanErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return &appconfig.Config.Vote.Nil, scanErr
	}
	return &reaction, nil
}
