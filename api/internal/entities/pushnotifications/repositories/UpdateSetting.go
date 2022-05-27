package repositories

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (s *repository) UpdateSetting(ctx context.Context, runnable utils.Runnable, userID int, req request.UpdatePushSetting) error {
	query, args, squirrelErr := squirrel.Update("push.settings").
		Set(req.ID, req.NewSetting).
		Where("user_id = ?", userID).
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
