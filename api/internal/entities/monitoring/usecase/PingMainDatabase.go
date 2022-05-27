package service

import (
	"context"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) PingMainDatabase(ctx context.Context) error {
	masterNode := svc.repo.MasterNode()
	dbErr := masterNode.Ping(ctx)

	if dbErr != nil {
		vlog.Errorf(ctx, "Error pinging DB: %v", dbErr)

		vlog.Infof(ctx, "SETTINGS:")
		vlog.Infof(ctx, "MAX size of the pool: %d", masterNode.Stat().MaxConns)
		vlog.Infof(ctx, "MAX duration (nanoseconds) since creation after which a connection will be automatically closed: %d", masterNode.Config().MaxConnLifetime)
		vlog.Infof(ctx, "MAX Connection idle time (nanoseconds): %d", masterNode.Config().MaxConnIdleTime)
		vlog.Infof(ctx, "MAX idle connections in the pool: %d", masterNode.Stat().IdleConns())

		vlog.Infof(ctx, "STATS:")
		vlog.Infof(ctx, "Num currently acquired connections in the pool: %d", masterNode.Stat().AcquiredConns())
		vlog.Infof(ctx, "Num cumulative count of successful acquires from the pool: %d", masterNode.Stat().AcquireCount())
		vlog.Infof(ctx, "Num idle connections in the pool: %d", masterNode.Stat().IdleConns())
		vlog.Infof(ctx, "Num acquires from the pool that were canceled by a context.: %d", masterNode.Stat().CanceledAcquireCount())
		return dbErr
	}

	pingErr := svc.repo.Redis().Ping(ctx)
	if pingErr != nil {
		return pingErr
	}

	return nil
}
