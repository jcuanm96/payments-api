package search

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	MasterNode() *pgxpool.Pool
	baserepo.BaseRepository
	GetChannelsByLink(ctx context.Context, query string, limit uint64, runnable utils.Runnable) ([]string, error)
	SearchUsersCalcQuery(ctx context.Context, userID int, filters []request.CustomFilterItem) (string, baserepo.SortDefinitionFunc, []string, error)
	SearchUsersExecQuery(ctx context.Context, query string, params []string) ([]response.User, error)
	GetOwnedByChannels(ctx context.Context, query string, limit uint64, runnable utils.Runnable) ([]string, error)
}
