package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
)

func (s *repository) IsFollowing(ctx context.Context, userID, goatUserID int) (bool, error) {
	query, args, squirrelErr := squirrel.Select(
		"id",
	).
		From("feed.follows").
		Where("user_id = ?", userID).
		Where("goat_user_id = ?", goatUserID).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return false, squirrelErr
	}

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		return false, queryErr
	}

	defer rows.Close()
	return rows.Next(), nil
}
