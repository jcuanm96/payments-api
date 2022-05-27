package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) RemoveFreeGroupChatCoCreator(ctx context.Context, runnable utils.Runnable, req request.RemoveFreeGroupChatCoCreator) error {
	query, args, squirrelErr := squirrel.Delete("product.free_group_creators").
		Where("sendbird_channel_id=?", req.ChannelID).
		Where("creator_user_id=?", req.UserID).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := runnable.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
