package repositories

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetPostConversation(ctx context.Context, postID int) (*response.Conversation, error) {
	query, args, squirrelErr := squirrel.Select(
		`messages.sendbird_channel_id`,
		`COALESCE(messages.conversation_start_ts, 0)`,
		`COALESCE(messages.conversation_end_ts, 0)`,
	).
		From(`feed.goat_chat_messages messages`).
		Join(`feed.posts posts ON messages.id = posts.goat_chat_msgs_id`).
		Where(`posts.id = ?`, postID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	rows, queryErr := s.MasterNode().Query(context.Background(), query, args...)
	if queryErr != nil {
		vlog.Errorf(ctx, "Error executing query: %s", query)
		return nil, queryErr
	}

	defer rows.Close()
	hasNext := rows.Next()

	if !hasNext {
		noRowsErr := fmt.Errorf("no rows returned for query: %s. post_id: %d", query, postID)
		vlog.Errorf(ctx, noRowsErr.Error())
		return nil, noRowsErr
	}

	conversation := response.Conversation{}
	scanErr := rows.Scan(
		&conversation.ChannelID,
		&conversation.StartTS,
		&conversation.EndTS,
	)
	if scanErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, nil))
		return nil, scanErr
	}

	return &conversation, nil
}
