package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) AddFreeGroupChatCoCreators(ctx context.Context, runnable utils.Runnable, req request.AddFreeGroupChatCoCreators) error {
	queryBuilder := squirrel.Insert("product.free_group_creators").
		Columns(
			"sendbird_channel_id",
			"creator_user_id",
		)

	for _, coCreator := range req.CoCreators {
		queryBuilder = queryBuilder.Values(
			req.ChannelID,
			coCreator,
		)
	}

	query, args, squirrelErr := queryBuilder.
		Suffix(`ON CONFLICT DO NOTHING`).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := runnable.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Errorf(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
