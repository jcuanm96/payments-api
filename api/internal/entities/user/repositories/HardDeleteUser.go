package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) HardDeleteUser(ctx context.Context, userID int) error {
	query, args, squirrelErr := squirrel.Delete("core.users").
		Where("id=?", userID).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}
	_, queryErr := s.MasterNode().Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
