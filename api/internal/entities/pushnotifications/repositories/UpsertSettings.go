package repositories

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (s *repository) UpsertSettings(ctx context.Context, runnable utils.Runnable, userID int, settings push.UpdateSettings) error {
	query, args, squirrelErr := squirrel.Insert("push.settings").
		Columns(
			"user_id",
			"pending_balance",
		).
		Values(
			userID,
			settings.PendingBalance,
		).
		Suffix(
			`ON CONFLICT (user_id)
		 	DO UPDATE SET 
			pending_balance = ?`,
			settings.PendingBalance,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := runnable.Exec(ctx, query, args...)
	if queryErr != nil {
		return errors.New(utils.SqlErrLogMsg(queryErr, query, args))
	}

	return nil
}
