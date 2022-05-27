package repositories

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func (s *repository) IsUserBlockedByPoster(ctx context.Context, userID int, postID int) (bool, error) {
	args := []interface{}{postID, userID}
	query := `
		SELECT user_id
		FROM core.user_blocks
		WHERE 
			blocked_user_id = $2 AND
			user_id IN (
				SELECT user_id
				FROM feed.posts
				WHERE id = $1
			)
	`

	row := s.MasterNode().QueryRow(ctx, query, args...)

	var id int
	scanErr := row.Scan(&id)

	if scanErr == pgx.ErrNoRows {
		return false, nil
	} else if scanErr != nil {
		return false, scanErr
	}

	return true, nil
}
