package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetFollowedGoats(ctx context.Context, userID int, req request.GetFollowedGoats) ([]response.User, error) {
	query := squirrel.Select(
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
		From("core.users users").
		Join("feed.follows follows ON users.id = follows.goat_user_id").
		Where("follows.user_id = ?", userID)

	if req.CursorGoatUserID > 0 {
		query = query.Where("users.id < ?", req.CursorGoatUserID)
	}

	queryString, args, queryErr := query.
		OrderBy("users.id DESC").
		Limit(uint64(req.Limit)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if queryErr != nil {
		return []response.User{}, queryErr
	}

	rows, sqlErr := s.MasterNode().Query(ctx, queryString, args...)
	if sqlErr != nil {
		return []response.User{}, sqlErr
	}

	defer rows.Close()
	goats := []response.User{}

	for rows.Next() {
		currUser := response.User{}
		scanErr := rows.Scan(
			&currUser.ID,
			&currUser.FirstName,
			&currUser.LastName,
			&currUser.Phonenumber,
			&currUser.CountryCode,
			&currUser.Email,
			&currUser.Username,
			&currUser.Type,
			&currUser.ProfileAvatar,
		)
		if scanErr != nil {
			vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, queryString, args))
			return nil, scanErr
		}

		goats = append(goats, currUser)
	}

	return goats, nil
}
