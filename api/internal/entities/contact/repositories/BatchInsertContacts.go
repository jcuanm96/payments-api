package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) BatchInsertContacts(ctx context.Context, currUserID int, userIDs []int, exec utils.Executable) error {
	queryBuilder := squirrel.Insert("core.users_contacts").
		Columns(
			"user_id",
			"contact_id",
		)

	for _, contactUserID := range userIDs {
		queryBuilder = queryBuilder.Values(
			contactUserID,
			currUserID,
		)
	}

	query, args, squirrelErr := queryBuilder.
		Suffix("ON CONFLICT DO NOTHING").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := exec.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
