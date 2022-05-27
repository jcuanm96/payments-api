package repositories

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/google/uuid"
)

func (s *repository) CensorUserData(ctx context.Context, exec utils.Executable, userID int) error {
	query, args, squirrelErr := squirrel.Update("core.users").
		Set("phone_number", uuid.New()).
		Set("email", uuid.New()).
		Set("first_name", "[deleted]").
		Set("last_name", "").
		Set("username", uuid.New()).
		Set("profile_avatar", nil).
		Set("deleted_at", time.Now()).
		Where("id = ?", userID).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := exec.Exec(ctx, query, args...)
	if queryErr != nil {
		return queryErr
	}

	return nil
}
