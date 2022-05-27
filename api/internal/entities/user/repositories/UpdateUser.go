package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) UpdateUser(ctx context.Context, runnable utils.Runnable, userID int, item request.UpdateUser) error {
	q := squirrel.Update("core.users")
	if item.FirstName != nil {
		q = q.Set("first_name", *item.FirstName)
	}
	if item.LastName != nil {
		q = q.Set("last_name", *item.LastName)
	}
	if item.Username != nil {
		q = q.Set("username", *item.Username)
	}
	if item.Email != nil {
		q = q.Set("email", *item.Email)
	}
	if item.ProfileAvatar != nil {
		q = q.Set("profile_avatar", *item.ProfileAvatar)
	}

	query, args, squirrelErr := q.Where("id=?", userID).
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
