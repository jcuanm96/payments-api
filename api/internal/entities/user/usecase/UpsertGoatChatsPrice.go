package service

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func (svc *usecase) UpsertGoatChatsPrice(ctx context.Context, tx pgx.Tx, priceInSmallestDenom int64, currency string, goatUserID int) error {
	return svc.repo.UpsertGoatChatsPrice(ctx, tx, priceInSmallestDenom, currency, goatUserID)
}
