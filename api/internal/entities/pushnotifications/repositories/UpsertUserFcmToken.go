package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
)

func (s *repository) UpsertUserFcmToken(ctx context.Context, userID int, token string) error {
	query, args, squirrelErr := squirrel.Insert("core.tokens").
		Columns(
			"user_id",
			"fcm_token",
		).
		Values(
			userID,
			token,
		).
		Suffix(`
		ON CONFLICT (user_id)
		DO UPDATE SET
			fcm_token = ?`,
			token,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := s.MasterNode().Exec(ctx, query, args...)
	if queryErr != nil {
		return queryErr
	}

	return nil
}
