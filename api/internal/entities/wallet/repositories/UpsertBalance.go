package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func UpsertBalance(ctx context.Context, runnable utils.Runnable, ledgerEntry wallet.LedgerEntry) error {
	query := `
		INSERT INTO wallet.balances(
			provider_user_id,
			available_balance,
			currency
		) VALUES (
			$1,
			$2,
			$3
		)
		ON CONFLICT (
			provider_user_id,
			currency
		) DO UPDATE
		SET available_balance = wallet.balances.available_balance + EXCLUDED.available_balance;
	`
	queryArgs := []interface{}{
		ledgerEntry.ProviderUserID,
		ledgerEntry.BalanceDelta,
		ledgerEntry.Currency,
	}

	_, queryErr := runnable.Exec(ctx, query, queryArgs...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, queryArgs))
		return queryErr
	}

	return nil
}

func (s *repository) UpsertBalance(ctx context.Context, runnable utils.Runnable, ledgerEntry wallet.LedgerEntry) error {
	return UpsertBalance(ctx, runnable, ledgerEntry)
}
