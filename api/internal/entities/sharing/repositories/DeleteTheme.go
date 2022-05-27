package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) DeleteTheme(ctx context.Context, themeID int) error {
	query, args, squirrelErr := squirrel.Delete("sharing.themes").
		Where("id=?", themeID).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, execErr := s.db.Exec(ctx, query, args...)
	if execErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(execErr, query, args))
		return execErr
	}

	return nil
}
