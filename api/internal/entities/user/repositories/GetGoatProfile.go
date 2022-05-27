package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (s *repository) GetGoatProfile(ctx context.Context, goatID, currUserID int) (*response.GetGoatProfile, error) {
	query := `
		WITH following AS (
			SELECT 
				id
			FROM feed.follows
			WHERE 
				user_id = $1 AND
				goat_user_id = $2
		),
		notifications AS (
			SELECT
				id
			FROM feed.post_notifications
			WHERE 
				user_id = $1 AND
				goat_user_id = $2
		)
		SELECT 
			users.id,
			first_name,
			last_name,
			phone_number,
			country_code,
			email,
			username,
			user_type,
			profile_avatar,
			COALESCE(following.id, 0),
			COALESCE(notifications.id, 0)
		FROM core.users users
		LEFT JOIN following ON TRUE
		LEFT JOIN notifications ON TRUE
		WHERE users.id = $2
 	`
	args := []interface{}{currUserID, goatID}
	row := s.MasterNode().QueryRow(ctx, query, args...)

	res := &response.GetGoatProfile{}
	var followCount int
	var notificationCount int
	scanErr := row.Scan(
		&res.User.ID,
		&res.User.FirstName,
		&res.User.LastName,
		&res.User.Phonenumber,
		&res.User.CountryCode,
		&res.User.Email,
		&res.User.Username,
		&res.User.Type,
		&res.User.ProfileAvatar,
		&followCount,
		&notificationCount,
	)

	if scanErr != nil {
		return nil, scanErr
	}

	res.IsFollowing = followCount > 0
	res.PostNotificationsEnabled = notificationCount > 0

	return res, nil
}
