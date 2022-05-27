package service

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func (svc *usecase) UseInviteCode(ctx context.Context, tx pgx.Tx, code string, userID int) error {
	return svc.repo.UseInviteCode(ctx, tx, code, userID)
}
