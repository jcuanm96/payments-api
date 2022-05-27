package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetMyInvites(ctx context.Context, userID int) ([]response.Invite, error) {
	query, args, squirrelErr := squirrel.Select(
		"codes.invite_code",

		"COALESCE(users.id, 0)",
		"COALESCE(users.first_name, '')",
		"COALESCE(users.last_name, '')",
		"COALESCE(users.phone_number, '')",
		"COALESCE(users.country_code, '')",
		"COALESCE(users.email, '')",
		"COALESCE(users.username, '')",
		"COALESCE(users.user_type, '')",
		"COALESCE(users.profile_avatar, '')",
	).
		From("core.goat_invite_codes codes").
		LeftJoin("core.users users ON users.id = codes.used_by").
		Where("codes.invited_by = ?", userID).
		OrderBy("codes.id DESC").
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
	invites := []response.Invite{}

	for rows.Next() {
		invite := response.Invite{}

		user := response.User{}
		scanErr := rows.Scan(
			&invite.Code,

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

		if user.ID != 0 {
			invite.User = &user
		}
		invites = append(invites, invite)
	}
	return invites, nil
}
