package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/jackc/pgx/v4"
)

const hardDeleteUserErr = "Something went wrong when trying to delete user."

func (svc *usecase) HardDeleteUser(ctx context.Context) error {
	user, userErr := svc.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			hardDeleteUserErr,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			hardDeleteUserErr,
			"GetCurrentUser returned nil for HardDeleteUser",
		)
	}

	// DB protocol
	tx, txErr := svc.repo.MasterNode().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			hardDeleteUserErr,
			fmt.Sprintf("Could not create transaction: %v", txErr),
		)
	}

	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	censorErr := svc.repo.CensorUserData(ctx, tx, user.ID)
	if censorErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			hardDeleteUserErr,
			fmt.Sprintf("Error censoring data in db for user %d: %v", user.ID, censorErr),
		)
	}

	clearFollowsErr := svc.repo.ClearFollows(ctx, tx, user.ID)
	if clearFollowsErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			hardDeleteUserErr,
			fmt.Sprintf("Error clearing follows from db for user %d: %v", user.ID, clearFollowsErr),
		)
	}

	clearFCMTokensErr := svc.repo.ClearFCMTokens(ctx, tx, user.ID)
	if clearFCMTokensErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			hardDeleteUserErr,
			fmt.Sprintf("Error clearing fcm tokens from db for user %d: %v", user.ID, clearFCMTokensErr),
		)
	}

	clearAuthTokensErr := svc.repo.ClearAuthTokens(ctx, tx, user.ID)
	if clearAuthTokensErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			hardDeleteUserErr,
			fmt.Sprintf("Error clearing auth tokens from db for user %d: %v", user.ID, clearAuthTokensErr),
		)
	}

	// SendBird Protocol
	deleteUserErr := svc.sendBirdClient.DeleteUser(user.ID)
	if deleteUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			hardDeleteUserErr,
			fmt.Sprintf("Error getting list of channels for user %d in HardDeleteUser: %v", user.ID, deleteUserErr),
		)
	}

	commit = true
	return nil
}
