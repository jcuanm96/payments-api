package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
)

func (s *repository) IsBlocked(ctx context.Context, currUserID int, userID int) (bool, error) {
	query, args, squirrelErr := squirrel.Select(
		"id",
	).
		From("core.user_blocks").
		Where("user_id = ?", currUserID).
		Where("blocked_user_id = ?", userID).
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
