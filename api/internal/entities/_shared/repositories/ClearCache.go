package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ClearCache(db *pgxpool.Pool, r vredis.Client, ctx context.Context, key, prefix string) error {
	query, args, squirrelErr := squirrel.Delete("cache").
		Where("key=?", key).
		Where("prefix=?", prefix).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := db.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}
	return nil
}
