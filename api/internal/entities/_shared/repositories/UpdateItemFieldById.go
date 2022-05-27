package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"

	"github.com/Masterminds/squirrel"
)

func UpdateItemFieldById(s BaseRepository, ctx context.Context, tableName string, id int, field string, value interface{}) error {
	{
		query, args, squirrelErr := squirrel.Update(tableName).
			Set(field, value).
			Where("id=?", id).
			PlaceholderFormat(squirrel.Dollar).ToSql()

		if squirrelErr != nil {
			return squirrelErr
		}

		_, queryErr := s.MasterNode().Exec(ctx, query, args...)
		if queryErr != nil {
			vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
			return queryErr
		}
	}

	return nil
}
