package repositories

import (
	"context"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetFeedPosts(ctx context.Context, userID int, params feed.GetFeedPostsParams) ([]response.FeedPost, error) {
	queryf := `
	SELECT 
		feed.id,
		feed.user_id,
		feed.post_created_at,
		feed.post_num_upvotes,
		feed.post_num_downvotes,
		feed.post_text_content,
		COALESCE(feed.post_link_suffix, ''),
		feed.goat_chat_messages_sendbird_id,
		feed.goat_chat_messages_convo_start_ts,
		feed.goat_chat_messages_convo_end_ts,
		feed.num_comments,
		feed.post_img_url,
		feed.post_img_width,
		feed.post_img_height,
		%s, -- reaction subquery

		users.id,
		users.first_name,
		users.last_name,
		users.phone_number,
		users.country_code,
		users.email,
		users.username,
		users.user_type,
		users.profile_avatar,

		COALESCE(users2.id, 0),
		COALESCE(users2.first_name, ''),
		COALESCE(users2.last_name, ''),
		COALESCE(users2.phone_number, ''),
		COALESCE(users2.country_code, ''),
		COALESCE(users2.email, ''),
		COALESCE(users2.username, ''),
		COALESCE(users2.user_type, ''),
		COALESCE(users2.profile_avatar, '')

	FROM feed.feed_posts_view feed

	%s -- reaction join
	JOIN core.users users ON feed.user_id = users.id
	LEFT JOIN core.users users2 ON feed.goat_chat_messages_customer_id = users2.id

	%s -- where clause

	ORDER BY feed.id DESC

	%s -- limit clause
	`

	var postQuery string
	var args []interface{}
	if params.PostID != nil { // GetFeedPostByID
		postQuery, args = formatGetFeedPostByID(queryf, userID, params)
	} else if params.LinkSuffix != nil { //GetFeedPostByLinkSuffix
		postQuery, args = formatGetFeedPostByLinkSuffix(queryf, userID, params)
	} else if params.GoatUserID != nil { //GetGoatFeedPosts
		postQuery, args = formatGetGoatFeedPosts(queryf, userID, params)
	} else if params.CursorPostID != nil { // GetUserFeedPosts
		postQuery, args = formatGetUserFeedPosts(queryf, userID, params)
	} else {
		return nil, fmt.Errorf("unknown params to GetFeedPosts")
	}

	rows, queryErr := s.MasterNode().Query(ctx, postQuery, args...)
	if queryErr != nil {
		vlog.Errorf(ctx, "Erroring query: %s", postQuery)
		return make([]response.FeedPost, 0), queryErr
	}

	defer rows.Close()
	posts := make([]response.FeedPost, 0)
	for rows.Next() {
		currFeedPost := response.FeedPost{}
		creator := response.User{}
		customer := &response.User{}
		var linkSuffix string
		scanErr := rows.Scan(
			&currFeedPost.ID,
			&currFeedPost.UserID,
			&currFeedPost.PostCreatedAt,
			&currFeedPost.PostNumUpvotes,
			&currFeedPost.PostNumDownvotes,
			&currFeedPost.PostTextContent,
			&linkSuffix,
			&currFeedPost.Conversation.ChannelID,
			&currFeedPost.Conversation.StartTS,
			&currFeedPost.Conversation.EndTS,
			&currFeedPost.NumComments,
			&currFeedPost.Image.URL,
			&currFeedPost.Image.Width,
			&currFeedPost.Image.Height,
			&currFeedPost.Reaction,

			&creator.ID,
			&creator.FirstName,
			&creator.LastName,
			&creator.Phonenumber,
			&creator.CountryCode,
			&creator.Email,
			&creator.Username,
			&creator.Type,
			&creator.ProfileAvatar,

			&customer.ID,
			&customer.FirstName,
			&customer.LastName,
			&customer.Phonenumber,
			&customer.CountryCode,
			&customer.Email,
			&customer.Username,
			&customer.Type,
			&customer.ProfileAvatar,
		)
		if scanErr != nil {
			vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, postQuery, args))
			return nil, scanErr
		}

		if customer.ID < 1 {
			customer = nil
		}

		currFeedPost.Creator = creator
		currFeedPost.Customer = customer

		redirectBaseURL := appconfig.Config.Gcloud.RedirectBaseURL
		currFeedPost.Link = fmt.Sprintf(constants.FEED_POST_BASE_URL_F, redirectBaseURL, linkSuffix)

		const length = constants.MAX_POST_PREVIEW_CONTENT_LENGTH
		if !params.IsTextContentFullLength && len(currFeedPost.PostTextContent) > length {
			currFeedPost.PostTextContent = currFeedPost.PostTextContent[:length]
		}

		posts = append(posts, currFeedPost)
	}
	return posts, nil
}

func formatGetFeedPostByID(queryf string, userID int, params feed.GetFeedPostsParams) (string, []interface{}) {
	reactionSelect := `
			COALESCE(
				(
					SELECT type
					FROM feed.post_reactions
					WHERE
						user_id = $1
						AND post_id = $2
						AND deleted_at IS NULL
				), 
				''
			)
		`
	reactionJoin := ``
	whereClause := ` 
			WHERE
				feed.id = $2 AND 
				feed.post_deleted_at IS NULL
		`
	limitString := ""
	args := []interface{}{userID, *params.PostID}
	postQuery := fmt.Sprintf(queryf, reactionSelect, reactionJoin, whereClause, limitString)
	return postQuery, args
}

func formatGetFeedPostByLinkSuffix(queryf string, userID int, params feed.GetFeedPostsParams) (string, []interface{}) {
	reactionSelect := `COALESCE(reacts.type, '')`
	reactionJoin := `
		LEFT JOIN feed.post_reactions reacts ON 
			feed.id = reacts.post_id AND 
			reacts.user_id = $1
	`
	whereClause := ` 
			WHERE
				feed.post_link_suffix = $2 AND 
				feed.post_deleted_at IS NULL
		`
	limitString := ""
	args := []interface{}{userID, *params.LinkSuffix}
	postQuery := fmt.Sprintf(queryf, reactionSelect, reactionJoin, whereClause, limitString)
	return postQuery, args
}

func formatGetGoatFeedPosts(queryf string, userID int, params feed.GetFeedPostsParams) (string, []interface{}) {
	reactionSelect := `COALESCE(reacts.type, '')`
	reactionJoin := `
		LEFT JOIN feed.post_reactions reacts ON 
			feed.id = reacts.post_id AND 
			reacts.user_id = $1
	`
	args := []interface{}{userID, *params.GoatUserID}
	whereClause := `
		WHERE 
			feed.user_id = $2 AND
			feed.post_deleted_at IS NULL	
		`
	if *params.CursorPostID >= 1 {
		whereClause += " AND (feed.id < $3)"
		args = append(args, *params.CursorPostID)
	}
	limitString := fmt.Sprintf("LIMIT %d", *params.Limit)
	postQuery := fmt.Sprintf(queryf, reactionSelect, reactionJoin, whereClause, limitString)
	return postQuery, args
}

func formatGetUserFeedPosts(queryf string, userID int, params feed.GetFeedPostsParams) (string, []interface{}) {
	reactionSelect := `COALESCE(reacts.type, '')`
	reactionJoin := `
			LEFT JOIN feed.post_reactions reacts ON 
				feed.id = reacts.post_id AND 
				reacts.user_id = $1
		`
	whereClause := `
		WHERE 
			(feed.user_id = $1 OR 
			feed.user_id IN (SELECT follows.goat_user_id FROM feed.follows AS follows WHERE follows.user_id = $1))
			AND feed.post_deleted_at IS NULL
		`
	args := []interface{}{userID}
	if *params.CursorPostID >= 1 {
		whereClause += " AND (feed.id < $2)"
		args = append(args, *params.CursorPostID)
	}
	limitString := fmt.Sprintf("LIMIT %d", *params.Limit)
	postQuery := fmt.Sprintf(queryf, reactionSelect, reactionJoin, whereClause, limitString)
	return postQuery, args
}
