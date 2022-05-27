package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	db  *pgxpool.Pool
	rdb vredis.Client
}

func (r *repository) MasterNode() *pgxpool.Pool {
	return r.db
}

func New(
	db *pgxpool.Pool,
	rdb vredis.Client,
) *repository {
	return &repository{
		db:  db,
		rdb: rdb,
	}
}

// Base cache repo functions
func (r *repository) ClearCache(ctx context.Context, key, prefix string) error {
	return baserepo.ClearCache(r.MasterNode(), r.rdb, ctx, key, prefix)
}
func (r *repository) GetFromCache(ctx context.Context, runnable utils.Runnable, key, prefix string, obj interface{}) error {
	return baserepo.GetFromCache(ctx, runnable, r.rdb, key, prefix, obj)
}
func (r *repository) SaveToCacheForSeconds(ctx context.Context, tx pgx.Tx, key, prefix string, seconds int, obj interface{}, userID int) error {
	return baserepo.SaveToCacheForSeconds(tx, r.rdb, ctx, key, prefix, seconds, obj, userID)
}
func (r *repository) SaveToCacheForHours(ctx context.Context, tx pgx.Tx, key, prefix string, hours int, obj interface{}, userID int) error {
	return baserepo.SaveToCacheForSeconds(tx, r.rdb, ctx, key, prefix, hours*60*60, obj, userID)
}

// Base repo functions
func (s *repository) CalcQueryWithSortAndOffset(ctx context.Context, q string, sorts []request.SortItem, pageSize, pageNumber int, sortDefinition baserepo.SortDefinitionFunc) (string, error) {
	return baserepo.CalcQueryWithSortAndOffset(ctx, q, sorts, pageSize, pageNumber, sortDefinition)
}

func (s *repository) GetCountByQuery(ctx context.Context, query string, params []string) (int, error) {
	return baserepo.GetCountByQuery(s, ctx, query, params)
}

func (s *repository) CalcPages(ctx context.Context, query string, params []string, reqPageNumber, reqPageSize, currentCount int) (response.Paging, error) {
	return baserepo.CalcPages(s, ctx, query, params, reqPageNumber, reqPageSize, currentCount)
}

func (s *repository) UpdateItemFieldById(ctx context.Context, tableName string, id int, field string, value interface{}) error {
	return baserepo.UpdateItemFieldById(s, ctx, tableName, id, field, value)
}

func (s *repository) FinishTx(ctx context.Context, tx pgx.Tx, commit *bool) error {
	return baserepo.FinishTx(ctx, tx, commit)
}
