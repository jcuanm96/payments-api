package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const blockUserDefaultErr = "Something went wrong when blocking user."

func (svc *usecase) BlockUser(ctx context.Context, blockUserID int) error {
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
			"Current user was nil in BlockUser",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			blockUserDefaultErr,
			fmt.Sprintf("Could not begin transaction: %v", txErr),
		)
	}
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	blockUserErr := svc.repo.BlockUser(ctx, tx, currUser.ID, blockUserID)
	if blockUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			blockUserDefaultErr,
			fmt.Sprintf("Error user %d blocking user %d: %v", currUser.ID, blockUserID, blockUserErr),
		)
	}

	_, sendbirdBlockUserErr := svc.sendBirdClient.BlockUser(currUser.ID, blockUserID)
	if sendbirdBlockUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			blockUserDefaultErr,
			fmt.Sprintf("Error user %d blocking user %d on Sendbird: %v", currUser.ID, blockUserID, sendbirdBlockUserErr),
		)
	}

	commit = true
	return nil
}
