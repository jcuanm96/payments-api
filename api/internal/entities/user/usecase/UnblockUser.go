package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errUnblockingUser = "Something went wrong when unblocking user."

func (svc *usecase) UnblockUser(ctx context.Context, unblockUserID int) error {
	currUser, userErr := svc.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if currUser == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"Current user was nil in UnblockUser",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUnblockingUser,
			fmt.Sprintf("Could not begin transaction: %v", txErr),
		)
	}

	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	unblockUserErr := svc.repo.UnblockUser(ctx, tx, currUser.ID, unblockUserID)
	if unblockUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUnblockingUser,
			fmt.Sprintf("Error user %d unblocking user %d: %v", currUser.ID, unblockUserID, unblockUserErr),
		)
	}

	sendbirdUnblockUserErr := svc.sendBirdClient.UnblockUser(currUser.ID, unblockUserID)
	if sendbirdUnblockUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUnblockingUser,
			fmt.Sprintf("Error user %d unblocking user %d on Sendbird: %v", currUser.ID, unblockUserID, sendbirdUnblockUserErr),
		)
	}

	commit = true
	return nil
}
