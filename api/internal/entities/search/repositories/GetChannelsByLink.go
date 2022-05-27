package repositories

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetChannelsByLink(ctx context.Context, query string, limit uint64, runnable utils.Runnable) ([]string, error) {
	likeClause := fmt.Sprintf("LOWER(link_suffix) LIKE '%s%s%s'", "%", query, "%")
	query, args, squirrelErr := squirrel.Select(
		"sendbird_channel_id",
	).
		From("product.paid_group_chats").
		Where(likeClause).
		Limit(limit).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	rows, sqlErr := runnable.Query(ctx, query, args...)
	if sqlErr != nil {
		return nil, sqlErr
	}

	defer rows.Close()
	channelIDs := []string{}

	for rows.Next() {
		var currChannelID string
		scanErr := rows.Scan(
			&currChannelID,
		)
		if scanErr != nil {
			vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
			return nil, scanErr
		}

		channelIDs = append(channelIDs, currChannelID)
	}

	return channelIDs, nil
}
