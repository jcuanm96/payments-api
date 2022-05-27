package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) UpdatePendingContacts(ctx context.Context, phone string, userID int, exec utils.Executable) error {
	query, args, squirrelErr := squirrel.Update("core.pending_contacts").
		Set("signed_up_user_id", userID).
		Where("phone_number=?", phone).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, execErr := exec.Exec(ctx, query, args...)
	if execErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(execErr, query, args))
		return execErr
	}

	return nil
}
