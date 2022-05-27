package service

import (
	"context"
	"fmt"
	"net/http"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) BatchAddUsersToContacts(ctx context.Context, req cloudtasks.AddUserContactsTask) error {
	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not begin transaction in BatchAddUsersToContacts. Error: %v", txErr),
		)
	}

	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	insertContactsErr := svc.repo.BatchInsertContacts(ctx, req.UserID, req.ContactUserIDs, tx)
	if insertContactsErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Something went wrong when batch inserting user contacts for user %d. Error: %v", req.UserID, insertContactsErr),
		)
	}

	user, getUserErr := svc.user.GetUserByID(ctx, req.UserID)
	if getUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Something went wrong when getting user %d. Error: %v", req.UserID, getUserErr),
		)
	}

	updateErr := svc.repo.UpdatePendingContacts(ctx, user.Phonenumber, user.ID, tx)
	if updateErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Something went wrong when updating pending contacts for user %d. Error: %v", req.UserID, getUserErr),
		)
	}

	commit = true

	return nil
}
