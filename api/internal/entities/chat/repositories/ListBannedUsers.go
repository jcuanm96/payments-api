package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) ListBannedUsers(ctx context.Context, runnable utils.Runnable, req request.ListBannedUsers) (*response.ListBannedUsers, error) {
	queryBuilder := squirrel.Select(
		"users.id",
		"users.first_name",
		"users.last_name",
		"users.phone_number",
		"users.country_code",
		"users.email",
		"users.username",
		"users.user_type",
		"users.profile_avatar",
	).
		From("core.banned_chat_users bans").
		Join("core.users users ON bans.banned_user_id = users.id").
		Where("bans.sendbird_channel_id = ?", req.ChannelID)

	if req.CursorID > 0 {
		queryBuilder = queryBuilder.Where("bans.banned_user_id < ?", req.CursorID)
	}

	query, args, squirrelErr := queryBuilder.
		OrderBy("bans.banned_user_id DESC").
		Limit(uint64(req.Limit)).
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

	users := []response.User{}
	for rows.Next() {
		user := response.User{}
		scanErr := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Phonenumber,
			&user.CountryCode,
			&user.Email,
			&user.Username,
			&user.Type,
			&user.ProfileAvatar,
		)
		if scanErr != nil {
			return nil, scanErr
		}

		users = append(users, user)
	}
	res := &response.ListBannedUsers{
		Users: users,
	}
	return res, nil
}
