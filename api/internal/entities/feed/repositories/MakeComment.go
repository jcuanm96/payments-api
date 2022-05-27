package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4/pgxpool"
)

func MakeCommentDB(ctx context.Context, req request.MakeComment, userID int, db *pgxpool.Pool) (*response.CommentMetadata, error) {
	query, args, squirrelErr := squirrel.Insert("feed.post_comments").
		Columns(
			"post_id",
			"user_id",
			"text_content",
		).
		Values(
			req.PostID,
			userID,
			req.Text,
		).
		Suffix(`
			RETURNING 
				id,
				user_id,
				post_id,
				text_content,
				created_at,
				updated_at
		`).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	row := db.QueryRow(ctx, query, args...)
	res := response.CommentMetadata{}
	scanErr := row.Scan(
		&res.ID,
		&res.UserID,
		&res.PostID,
		&res.TextContent,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if scanErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return nil, scanErr
	}

	return &res, nil
}

func (s *repository) MakeComment(ctx context.Context, req request.MakeComment, userID int) (*response.CommentMetadata, error) {
	return MakeCommentDB(ctx, req, userID, s.MasterNode())
}
