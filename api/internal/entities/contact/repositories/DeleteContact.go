package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) DeleteContact(ctx context.Context, userID, contactID int) error {
	query, args, squirrelErr := squirrel.Delete("core.users_contacts").
		Where("user_id=?", userID).Where("contact_id=?", contactID).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(err, query, args))
		return err
	}

	return nil
}
