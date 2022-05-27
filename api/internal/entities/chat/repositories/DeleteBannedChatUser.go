package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) DeleteBannedChatUser(ctx context.Context, bannedUserID int, channelID string, runnable utils.Runnable) error {
	query, args, squirrelErr := squirrel.Delete("core.banned_chat_users").
		Where("banned_user_id = ?", bannedUserID).
		Where("sendbird_channel_id = ?", channelID).
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
