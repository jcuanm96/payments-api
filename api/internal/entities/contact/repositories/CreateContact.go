package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateContactDB(ctx context.Context, userID int, contactID int, db *pgxpool.Pool) error {
	query, args, squirrelErr := squirrel.Insert("core.users_contacts").
		Columns(
			"user_id",
			"contact_id",
		).
		Values(
			userID,
			contactID,
		).
		Suffix(`ON conflict (user_id, contact_id) do nothing`).
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

func (s *repository) CreateContact(ctx context.Context, userID int, contactID int) error {
	return CreateContactDB(ctx, userID, contactID, s.MasterNode())
}
