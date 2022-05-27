package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) InsertBannedChatUser(ctx context.Context, bannedUserID int, userID int, channelID string, runnable utils.Runnable) error {
	query, args, squirrelErr := squirrel.Insert("core.banned_chat_users").
		Columns(
			"banned_user_id",
			"user_id",
			"sendbird_channel_id",
		).
		Values(
			bannedUserID,
			userID,
			channelID,
		).
		Suffix(`ON conflict (banned_user_id, sendbird_channel_id) do nothing`).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		vlog.Errorf(ctx, "Error constructing banned chat user insert query: %s", squirrelErr.Error())
		return squirrelErr
	}
	_, queryErr := runnable.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
