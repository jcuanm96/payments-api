package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4/pgxpool"
)

func FollowDB(ctx context.Context, userID, goatUserID int, db *pgxpool.Pool) error {
	query, args, squirrelErr := squirrel.Insert("feed.follows").
		Columns(
			"user_id",
			"goat_user_id",
		).
		Values(
			userID,
			goatUserID,
		).
		Suffix(`ON conflict (user_id, goat_user_id) do nothing`).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := db.Exec(ctx, query, args...)
	if queryErr != nil {
		if utils.DuplicateError(queryErr) {
			return constants.ErrAlreadyExists
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}

func (s *repository) Follow(ctx context.Context, userID int, goatUserID int) error {
	return FollowDB(ctx, userID, goatUserID, s.MasterNode())
}
