package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetBioData(ctx context.Context, runnable utils.Runnable, username string, userID int) (*response.GetBioData, error) {
	res := response.GetBioData{}

	selectQuery := `
		SELECT 
			u.id,
			u.first_name, 
			u.last_name, 
			data.text_content,
			u.profile_avatar,
			(SELECT COUNT(*) FROM feed.goat_chat_messages WHERE provider_user_id = $1) AS num_goat_chats,
			(SELECT COUNT(*) FROM subscription.user_subscriptions WHERE goat_user_id = $1) AS num_subscribers,
			(SELECT COUNT(*) FROM feed.feed_posts_view WHERE user_id = $1) AS num_feed_posts
		FROM core.users AS u
		LEFT JOIN core.user_bio data ON u.id = data.user_id
	  	WHERE
			lower(username) = $2
			AND deleted_at IS NULL;
	`

	args := []interface{}{userID, username}

	rows, err := runnable.Query(ctx, selectQuery, args...)
	if err != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(err, selectQuery, args))
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(
			&res.UserID,
			&res.FirstName,
			&res.LastName,
			&res.TextContent,
			&res.ProfileAvatar,
			&res.NumGoatChats,
			&res.NumSubscribers,
			&res.NumFeedPosts,
		)
		if err != nil {
			vlog.Error(ctx, utils.SqlErrLogMsg(err, selectQuery, args))
			return nil, err
		}
	}

	return &res, nil
}
