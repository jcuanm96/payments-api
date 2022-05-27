package repositories

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func (r *repository) MasterNode() *pgxpool.Pool {
	return r.db
}

func New(
	db *pgxpool.Pool,
) *repository {
	return &repository{
		db: db,
	}
}
