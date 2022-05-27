package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4/pgxpool"
)

func UpsertBioDB(ctx context.Context, userID int, bio string, db *pgxpool.Pool) error {
	query, args, squirrelErr := squirrel.Insert("core.user_bio").
		Columns(
			"user_id",
			"text_content",
		).
		Values(
			userID,
			bio,
		).
		Suffix(`ON conflict (user_id)
			do update SET 
			text_content = ?`,
			bio,
		).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}
	_, queryErr := db.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}

func (s *repository) UpsertBio(ctx context.Context, userID int, bio string) error {
	return UpsertBioDB(ctx, userID, bio, s.MasterNode())
}
