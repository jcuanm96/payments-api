package repositories

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func GetFromCache(ctx context.Context, runnable utils.Runnable, r vredis.Client, key, prefix string, obj interface{}) error {
	res := ""
	query, args, squirrelErr := squirrel.Select(
		"model",
	).From("cache").
		Where("key = ?", key).
		Where("prefix = ?", prefix).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}
	row := runnable.QueryRow(ctx, query, args...)
	scanErr := row.Scan(
		&res,
	)
	if scanErr != nil {
		if scanErr != pgx.ErrNoRows {
			vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
			obj = nil
			return scanErr
		}
	} else {
		unmarshalErr := json.Unmarshal([]byte(res), &obj)
		if unmarshalErr != nil {
			vlog.Error(ctx, unmarshalErr.Error())
			return unmarshalErr
		}
	}
	return nil
}
