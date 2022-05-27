package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) DeletePaidGroupSubscription(ctx context.Context, runnable utils.Runnable, channelID string, userID int) error {
	query, args, squirrelErr := squirrel.Delete("subscription.paid_group_chat_subscriptions").
		Where("sendbird_channel_id = ?", channelID).
		Where("user_id = ?", userID).
		PlaceholderFormat(squirrel.Dollar).ToSql()

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
