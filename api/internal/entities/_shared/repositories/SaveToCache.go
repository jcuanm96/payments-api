package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Masterminds/squirrel"
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func SaveToCache(tx pgx.Tx, r vredis.Client, ctx context.Context, key, prefix string, obj interface{}, dur time.Duration, userID int) error {
	expDate := time.Now().Add(dur).Unix()
	json, _ := json.Marshal(obj)
	strJson := string(json)
	query, args, sqlErr := squirrel.Insert("cache").
		Columns(
			"key",
			"prefix",
			"expired_at",
			"model",
			"user_id",
		).
		Values(
			key,
			prefix,
			expDate,
			strJson,
			userID,
		).
		Suffix(`ON conflict (key,prefix)
			do update SET 
			expired_at = ?,
			model = ?,
			user_id = ?`,
			expDate, strJson, userID).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if sqlErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(sqlErr, query, args))
		return sqlErr
	}
	_, queryErr := tx.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}
	return nil
}
