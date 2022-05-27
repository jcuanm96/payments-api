package repositories

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetOwnedByChannels(ctx context.Context, query string, limit uint64, runnable utils.Runnable) ([]string, error) {
	likeClausef := `
		LOWER(goat.first_name) LIKE '%s%s%s' OR
		LOWER(goat.last_name) LIKE '%s%s%s' OR
		LOWER(goat.username) LIKE '%s%s%s' OR 
		LOWER(pgc.link_suffix) LIKE '%s%s%s'
	`

	likeClause := fmt.Sprintf(
		likeClausef,
		"%", query, "%",
		"%", query, "%",
		"%", query, "%",
		"%", query, "%",
	)
	query, args, squirrelErr := squirrel.Select(
		"pgc.sendbird_channel_id",
	).
		From("product.paid_group_chats pgc").
		Join("core.users goat ON goat.id = pgc.goat_user_id").
		Where(likeClause).
		OrderBy(`
			goat.first_name,
			goat.last_name,
			goat.username,
			pgc.link_suffix
		`).
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
