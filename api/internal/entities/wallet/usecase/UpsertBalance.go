package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (svc *usecase) UpsertBalance(ctx context.Context, runnable utils.Runnable, ledgerEntry wallet.LedgerEntry) error {
	return svc.repo.UpsertBalance(ctx, runnable, ledgerEntry)
}
