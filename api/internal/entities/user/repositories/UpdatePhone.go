package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) UpdatePhone(ctx context.Context, runnable utils.Runnable, userID int, phone string) error {
	phoneQuery, phoneArgs, phoneErr := squirrel.Update("core.users").
		Set("phone_number", phone).
		Where("id=?", userID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if phoneErr != nil {
		vlog.Errorf(ctx, "Something went wrong when using squirrel to create the SQL string to update a user's phone number.  Err: %s\n", phoneErr.Error())
		return phoneErr
	}

	_, execQueryErr := runnable.Exec(ctx, phoneQuery, phoneArgs...)
	if execQueryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(execQueryErr, phoneQuery, phoneArgs))
		return execQueryErr
	}

	return nil
}
