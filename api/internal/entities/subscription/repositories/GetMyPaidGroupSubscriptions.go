package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetMyPaidGroupSubscriptions(ctx context.Context, userID int, cursorID int, limit uint64) ([]response.PaidGroupChatSubscription, error) {
	queryBuilder := squirrel.Select(
		"users.id",
		"users.first_name",
		"users.last_name",
		"users.username",
		"users.user_type",
		"users.profile_avatar",

		"subscriptions.id",
		"subscriptions.current_period_end",
		"subscriptions.sendbird_channel_id",
		"subscriptions.is_renewing",
	).
		From("core.users users").
		Join("subscription.paid_group_chat_subscriptions subscriptions ON subscriptions.goat_user_id = users.id").
		Where("subscriptions.user_id = ?", userID).
		Where("subscriptions.current_period_end > now()")

	if cursorID > 0 {
		queryBuilder = queryBuilder.Where("subscriptions.id < ?", cursorID)
	}

	query, args, squirrelErr := queryBuilder.
		OrderBy("subscriptions.id DESC").
		Limit(limit).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}
	defer rows.Close()

	subscriptions := []response.PaidGroupChatSubscription{}
	for rows.Next() {
		subscription := response.PaidGroupChatSubscription{}
		scanErr := rows.Scan(
			&subscription.GoatUser.ID,
			&subscription.GoatUser.FirstName,
			&subscription.GoatUser.LastName,
			&subscription.GoatUser.Username,
			&subscription.GoatUser.Type,
			&subscription.GoatUser.ProfileAvatar,

			&subscription.ID,
			&subscription.CurrentPeriodEnd,
			&subscription.ChannelID,
			&subscription.IsRenewing,
		)

		if scanErr != nil {
			return nil, scanErr
		}

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}
