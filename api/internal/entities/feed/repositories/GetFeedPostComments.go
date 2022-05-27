package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetFeedPostComments(ctx context.Context, req request.GetFeedPostComments, runnable utils.Runnable) ([]response.Comment, error) {
	args := []interface{}{req.PostID}
	query := squirrel.Select(
		// comments object
		"comments.id",
		"comments.user_id",
		"comments.post_id",
		"comments.text_content",
		"comments.created_at",
		"comments.updated_at",

		// user object
		"users.id",
		"users.uuid",
		"users.stripe_id",
		"users.first_name",
		"users.last_name",
		"users.phone_number",
		"users.country_code",
		"users.email",
		"users.username",
		"users.user_type",
		"users.profile_avatar",
		"users.created_at",
		"users.updated_at",
		"users.deleted_at",
		"users.stripe_account_id",
	).
		From("feed.post_comments comments").
		Join("core.users users ON comments.user_id = users.id").
		Where("comments.post_id = ?").
		Where("comments.deleted_at IS NULL")

	if req.LastUserCommentID >= 1 {
		query = query.Where("comments.id < ?")
		args = append(args, req.LastUserCommentID)
	}

	queryString, p, queryErr := query.
		OrderBy("comments.id DESC").
		Limit(uint64(req.Limit)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if queryErr != nil {
		return make([]response.Comment, 0), queryErr
	}

	rows, sqlErr := runnable.Query(ctx, queryString, args...)
	if sqlErr != nil {
		return make([]response.Comment, 0), sqlErr
	}

	defer rows.Close()
	comments := make([]response.Comment, 0)
	currComment := response.CommentMetadata{}
	currUser := response.User{}

	for rows.Next() {
		scanErr := rows.Scan(
			&currComment.ID,
			&currComment.UserID,
			&currComment.PostID,
			&currComment.TextContent,
			&currComment.CreatedAt,
			&currComment.UpdatedAt,
			&currUser.ID,
			&currUser.UUID,
			&currUser.StripeID,
			&currUser.FirstName,
			&currUser.LastName,
			&currUser.Phonenumber,
			&currUser.CountryCode,
			&currUser.Email,
			&currUser.Username,
			&currUser.Type,
			&currUser.ProfileAvatar,
			&currUser.CreatedAt,
			&currUser.UpdatedAt,
			&currUser.DeletedAt,
			&currUser.StripeAccountID,
		)
		if scanErr != nil {
			vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, queryString, p))
			return nil, scanErr
		}

		commentPayload := response.Comment{
			CommentMetadata: currComment,
			User:            currUser,
		}
		comments = append(comments, commentPayload)
	}

	return comments, nil
}
