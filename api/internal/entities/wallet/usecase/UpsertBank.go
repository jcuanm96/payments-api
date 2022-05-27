package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const upsertBankDefaultErr = "There was a problem saving your bank information."

func (svc *usecase) UpsertBank(ctx context.Context, req request.UpsertBank) error {
	currentUser, userErr := svc.user.GetCurrentUser(ctx)

	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			upsertBankDefaultErr,
			fmt.Sprintf("Error occurred when getting the user to update bank info. Err: %v", userErr),
		)
	} else if currentUser == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			upsertBankDefaultErr,
		)
	} else if currentUser.Type != "GOAT" {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can update their bank info. Interested?  Sign up to be a creator now!",
			"Only creators can update their bank info.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when updating bank info. Please try again.",
			fmt.Sprintf("Error starting tx: %v", txErr),
		)
	}

	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	upsertBankInfoErr := svc.repo.UpsertBank(ctx, tx, currentUser.ID, req)
	if upsertBankInfoErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when updating bank info. Please try again.",
			fmt.Sprintf("Error occurred when updating bank info for user %d. Err: %v", currentUser.ID, upsertBankInfoErr),
		)
	}

	upsertBillingAddressErr := svc.repo.UpsertBillingAddress(ctx, tx, currentUser.ID, req.BillingAddress)
	if upsertBillingAddressErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when updating billing address. Please try again.",
			fmt.Sprintf("Error occurred when updating billing address for user %d. Err: %v", currentUser.ID, upsertBillingAddressErr),
		)
	}

	commit = true
	return nil
}
