package webhooks

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	MasterNode() *pgxpool.Pool
}
