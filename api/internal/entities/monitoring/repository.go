package monitoring

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	MasterNode() *pgxpool.Pool
	Redis() vredis.Client
	GetDashboard(ctx context.Context) (*response.GetDashboard, error)
}
