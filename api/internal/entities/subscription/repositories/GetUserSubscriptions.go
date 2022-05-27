package repositories

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetUserSubscriptions(ctx context.Context, userID int, cursorID int, limit uint64) ([]response.UserSubscription, *int, error) {
	queryf := `
		WITH subscriptions AS (
			%s
		)
		SELECT
			users.id,
			first_name,
			last_name,
			username,
			user_type,
			profile_avatar,

			subscriptions.id,
			subscriptions.current_period_end,
			tier_id,
			subscriptions.is_renewing

		FROM core.users users
		JOIN subscriptions
			ON subscriptions.goat_user_id = users.id
		ORDER BY subscriptions.id DESC
	`

	withClause := squirrel.
		Select(
			"id",
			"user_id",
			"goat_user_id",
			"current_period_end",
			"tier_id",
			"is_renewing",
		).
		From("subscription.user_subscriptions").
		Where("user_id = ?", userID).
		Where("current_period_end > now()")

	if cursorID > 0 {
		withClause = withClause.Where("id < ?", cursorID)
	}

	withClauseStr, args, squirrelErr := withClause.
		OrderBy("id DESC").
		Limit(limit).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		vlog.Errorf(ctx, "Squirrel error building WITH clause: %s", squirrelErr.Error())
		return nil, nil, squirrelErr
	}

	query := fmt.Sprintf(queryf, withClauseStr)

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, nil, queryErr
	}

	subscriptions := []response.UserSubscription{}

	defer rows.Close()
	lastID := cursorID
	for rows.Next() {
		subscription := response.UserSubscription{}
		scanErr := rows.Scan(
			&subscription.GoatUser.ID,
			&subscription.GoatUser.FirstName,
			&subscription.GoatUser.LastName,
			&subscription.GoatUser.Username,
			&subscription.GoatUser.Type,
			&subscription.GoatUser.ProfileAvatar,

			&subscription.ID,
			&subscription.CurrentPeriodEnd,
			&subscription.TierID,
			&subscription.IsRenewing,
		)

		if subscription.ID < lastID || lastID == 0 {
			lastID = subscription.ID
		}

		if scanErr != nil {
			vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
			return nil, nil, scanErr
		}

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, &lastID, nil
}
