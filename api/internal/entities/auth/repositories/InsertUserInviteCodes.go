package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) InsertUserInviteCodes(ctx context.Context, exec utils.Executable, userID int, codes []string) error {
	queryBuilder := squirrel.Insert("core.goat_invite_codes").
		Columns(
			"invite_code",
			"invited_by",
		)
	for _, code := range codes {
		queryBuilder = queryBuilder.Values(
			code,
			userID,
		)
	}

	query, args, squirrelErr := queryBuilder.
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
