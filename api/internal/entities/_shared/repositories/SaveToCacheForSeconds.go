package repositories

import (
	"context"
	"time"

	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	"github.com/jackc/pgx/v4"
)

func SaveToCacheForSeconds(tx pgx.Tx, r vredis.Client, ctx context.Context, key, prefix string, seconds int, obj interface{}, userID int) error {
	return SaveToCache(tx, r, ctx, key, prefix, obj, time.Duration(seconds)*time.Second, userID)
}
