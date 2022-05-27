package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (s *repository) InsertMessageLink(ctx context.Context, runnable utils.Runnable, linkSuffix, messageID, channelID string) error {
	query, args, squirrelErr := squirrel.Insert("sharing.message_links").
		Columns(
			"link_suffix",
			"message_id",
			"sendbird_channel_id",
		).
		Values(
			linkSuffix,
			messageID,
			channelID,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := runnable.Exec(ctx, query, args...)
	if queryErr != nil {
		return queryErr
	}

	return nil
}
