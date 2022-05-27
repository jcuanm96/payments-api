package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	wallet "github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) MarkProviderAsPaid(ctx context.Context, req request.MarkProviderAsPaid) error {
	amountPaidNegative := req.AmountPaid * -1
	if amountPaidNegative > 0 {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Please enter a positive value for amountPaid.",
			"Please enter a positive value for amountPaid.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong marking provider as paid.",
			fmt.Sprintf("Error starting tx in MarkProviderAsPaid: %v", txErr),
		)
	}

	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	updateBalanceErr := svc.repo.UpdateBalancePayout(ctx, tx, req.ProviderID, amountPaidNegative, req.PayPeriodID)
	if updateBalanceErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong marking provider as paid.",
			fmt.Sprintf("Error updating balance: %v", updateBalanceErr),
		)
	}

	ledgerEntry := wallet.LedgerEntry{
		ProviderUserID: req.ProviderID,
		SourceType:     constants.TXN_HISTORY_PAYOUT,
		Amount:         amountPaidNegative,
		CreatedTS:      0,
		BalanceDelta:   amountPaidNegative,
		Currency:       constants.DEFAULT_CURRENCY,
		PayPeriodID:    req.PayPeriodID,
	}
	insertLedgerErr := svc.repo.InsertLedgerTransaction(ctx, tx, ledgerEntry)
	if insertLedgerErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong marking provider as paid.",
			fmt.Sprintf("Error inserting ledger entry: %v", insertLedgerErr),
		)
	}

	commit = true
	return nil
}
