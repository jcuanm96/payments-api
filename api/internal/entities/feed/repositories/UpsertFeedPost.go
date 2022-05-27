package repositories

import (
	"context"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/jackc/pgx/v4/pgxpool"
)

func UpsertFeedPostDB(ctx context.Context, req request.MakeFeedPost, userID int, image response.PostImage, linkSuffix string, db *pgxpool.Pool) (int, error) {
	// The following query does the following:
	// 1. Attempts to insert text_content, user_id, goat_chat_msgs_id into the feed.posts table
	// 2. If feed.posts insert successful and img_url is valid, insert into feed.post_images
	insertQuery := `
		WITH user_post_data(
			user_id, 
			text_content, 
			img_url, 
			img_width,
			img_height,
			goat_chat_msgs_id,
			link_suffix
		) AS (
			VALUES (
				CAST($1 AS INT), 
				$2, 
				$3,
				CAST($4 AS INT),
				CAST($5 AS INT),
				CASE  
					WHEN CAST($6 AS INT) = 0 THEN NULL 
					ELSE CAST($6 AS INT) 
				END,
				$7
			)
		),
		feed_posts_ins AS (
			INSERT INTO feed.posts(user_id, text_content, goat_chat_msgs_id, link_suffix)
			SELECT 
				user_post_data.user_id, 
				user_post_data.text_content,
				user_post_data.goat_chat_msgs_id,
				user_post_data.link_suffix
			FROM user_post_data
			RETURNING user_id, id
		),
		feed_post_images_ins AS (
			INSERT INTO feed.post_images(post_id, img_url, img_width, img_height)
			SELECT 
				feed_posts_ins.id, 
				user_post_data.img_url,
				user_post_data.img_width,
				user_post_data.img_height
			FROM user_post_data
			JOIN feed_posts_ins USING (user_id)
			WHERE user_post_data.img_url <> ''
		)
		SELECT id FROM feed_posts_ins
	`

	args := []interface{}{userID, req.TextContent, image.URL, image.Width, image.Height, req.GoatChatMessagesID, linkSuffix}
	row, queryErr := db.Query(ctx, insertQuery, args...)

	if queryErr != nil {
		return -1, queryErr
	}

	defer row.Close()
	hasNext := row.Next()
	if !hasNext {
		return -1, fmt.Errorf("no rows returned, expected new post ID. query: %s", insertQuery)
	}
	var id int
	scanErr := row.Scan(&id)
	if scanErr != nil {
		return -1, scanErr
	}

	return id, nil
}

func (s *repository) UpsertFeedPost(ctx context.Context, req request.MakeFeedPost, userID int, image response.PostImage, linkSuffix string) (int, error) {
	return UpsertFeedPostDB(ctx, req, userID, image, linkSuffix, s.MasterNode())
}
