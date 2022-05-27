package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func InsertLedgerTransactionTx(ctx context.Context, runnable utils.Runnable, ledgerEntry wallet.LedgerEntry) error {
	query := `
	INSERT INTO wallet.ledger (
		provider_user_id,
		customer_user_id,
		stripe_transaction_id,
		source_type,
		stripe_created_ts,
		stripe_fee,
		vama_fee,
		amount,
		currency,
		pay_period_id,
		version
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11
	)`

	customerUserID := &ledgerEntry.CustomerUserID
	if ledgerEntry.CustomerUserID <= 0 && ledgerEntry.CustomerUserID != constants.VAMA_USER_ID {
		customerUserID = nil
	}
	args := []interface{}{
		ledgerEntry.ProviderUserID,
		customerUserID,
		ledgerEntry.StripeTransactionID,
		ledgerEntry.SourceType,
		ledgerEntry.CreatedTS,
		ledgerEntry.StripeFee,
		ledgerEntry.VamaFee,
		ledgerEntry.Amount,
		ledgerEntry.Currency,
		ledgerEntry.PayPeriodID,
		constants.CURR_PAYMENTS_VERSION,
	}

	_, queryErr := runnable.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}

func (s *repository) InsertLedgerTransaction(ctx context.Context, runnable utils.Runnable, ledgerEntry wallet.LedgerEntry) error {
	return InsertLedgerTransactionTx(ctx, runnable, ledgerEntry)
}
