package repositories

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) UpdateGoatChat(ctx context.Context, req request.EndGoatChat) (*chat.UpdateGoatChatResult, error) {
	query, args, sqlErr := squirrel.Update("feed.goat_chat_messages").
		Set("conversation_start_ts", req.StartTS).
		Where(`
			id = (
				SELECT MAX(id)
				FROM feed.goat_chat_messages
				WHERE sendbird_channel_id = ?)
			`, req.SendBirdChannelID).
		Suffix(`
			RETURNING 
				id, 
				is_public
		`).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if sqlErr != nil {
		vlog.Errorf(ctx, "SQL err: %s\n", sqlErr.Error())
		return nil, sqlErr
	}

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}

	defer rows.Close()
	hasNext := rows.Next()

	if !hasNext {
		return nil, fmt.Errorf("no rows returned for query: %s. sendbird_channel_id: %s", query, req.SendBirdChannelID)
	}

	result := chat.UpdateGoatChatResult{}
	scanErr := rows.Scan(
		&result.GoatChatMessagesID,
		&result.IsPublic,
	)

	if scanErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, []interface{}{req.SendBirdChannelID}))
		return nil, scanErr
	}

	return &result, nil
}
