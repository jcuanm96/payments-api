package repositories

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
)

type repository struct {
	db *pgxpool.Pool
}

func (r *repository) MasterNode() *pgxpool.Pool {
	return r.db
}

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

func New(
	db *pgxpool.Pool,
) *repository {
	return &repository{
		db: db,
	}
}
