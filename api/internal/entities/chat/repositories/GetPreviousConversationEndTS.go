package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetPreviousConversationEndTS(ctx context.Context, channelID string) (int64, error) {
	query, args, selectErr := squirrel.Select(
		"COALESCE(conversation_end_ts, 0)",
	).
		From("feed.goat_chat_messages").
		Where(`
		id = (
			SELECT MAX(id)
			FROM feed.goat_chat_messages
			WHERE sendbird_channel_id = ?
			AND id NOT IN ( SELECT MAX(id) from feed.goat_chat_messages WHERE sendbird_channel_id = ?)
		)
		`, channelID, channelID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if selectErr != nil {
		return -1, selectErr
	}

	rows, queryErr := s.MasterNode().Query(context.Background(), query, args...)
	if queryErr != nil {
		return -1, queryErr
	}

	defer rows.Close()
	hasNext := rows.Next()

	if !hasNext {
		return 0, nil
	}

	var conversationEndTS int64
	scanErr := rows.Scan(
		&conversationEndTS,
	)

	if scanErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return -1, scanErr
	}
	return conversationEndTS, nil
}
