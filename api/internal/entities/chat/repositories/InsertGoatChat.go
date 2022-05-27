package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) InsertGoatChat(ctx context.Context, customerUserID int, providerUserID int, req request.StartGoatChat, runnable utils.Runnable) error {
	query, args, queryErr := squirrel.Insert("feed.goat_chat_messages").
		Columns(
			"sendbird_channel_id",
			"is_public",
			"customer_user_id",
			"provider_user_id",
		).
		Values(
			req.SendBirdChannelID,
			req.IsPublic,
			customerUserID,
			providerUserID,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if queryErr != nil {
		vlog.Errorf(ctx, "Error constructing creator chat insert query: %s", queryErr.Error())
		return queryErr
	}
	_, err := runnable.Exec(ctx, query, args...)
	if err != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(err, query, args))
		return err
	}

	return nil
}
