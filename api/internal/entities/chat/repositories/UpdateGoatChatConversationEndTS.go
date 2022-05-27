package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) UpdateGoatChatConversationEndTS(ctx context.Context, lastMessageCreatedAt int64, sendbirdChannelID string, runnable utils.Runnable) error {
	query, args, sqlErr := squirrel.Update("feed.goat_chat_messages").
		Set("conversation_end_ts", lastMessageCreatedAt).
		Where(`
			id = (
				SELECT MAX(id)
				FROM feed.goat_chat_messages
				WHERE sendbird_channel_id = ?)
			`, sendbirdChannelID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if sqlErr != nil {
		vlog.Errorf(ctx, "SQL err: %s\n", sqlErr.Error())
		return sqlErr
	}

	_, execErr := runnable.Exec(ctx, query, args...)
	if execErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(execErr, query, args))
		return execErr
	}

	return nil
}
