package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SortDefinitionFunc func(sort request.SortItem) string

type BaseRepository interface {
	MasterNode() *pgxpool.Pool
	CalcQueryWithSortAndOffset(ctx context.Context, q string, sorts []request.SortItem, pageSize, pageNumber int, sortDefinition SortDefinitionFunc) (string, error)
	CalcPages(ctx context.Context, query string, params []string, reqPageNumber, reqPageSize, currentCount int) (response.Paging, error)
	GetCountByQuery(ctx context.Context, query string, params []string) (int, error)
	UpdateItemFieldById(ctx context.Context, tableName string, id int, field string, value interface{}) error
	FinishTx(ctx context.Context, tx pgx.Tx, commit *bool) error
}

type BaseCacheRepository interface {
	ClearCache(ctx context.Context, key, prefix string) error
	SaveToCacheForSeconds(ctx context.Context, tx pgx.Tx, key string, prefix string, seconds int, obj interface{}, userID int) error
	SaveToCacheForHours(ctx context.Context, tx pgx.Tx, key string, prefix string, hours int, obj interface{}, userID int) error
	GetFromCache(ctx context.Context, runnable utils.Runnable, key string, prefix string, obj interface{}) error
}
