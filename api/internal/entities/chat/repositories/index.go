package repositories

import (
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	db  *pgxpool.Pool
	rdb vredis.Client
}

func (r *repository) MasterNode() *pgxpool.Pool {
	return r.db
}

func (r *repository) Redis() vredis.Client {
	return r.rdb
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
