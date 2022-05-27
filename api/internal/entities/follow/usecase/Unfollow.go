package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const unfollowDefaultErr = "Something went wrong when attempting to unfollow."

func (svc *usecase) Unfollow(ctx context.Context, userToUnfollowID int) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			unfollowDefaultErr,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			unfollowDefaultErr,
			"user returned nil in Unfollow",
		)
	}

	deleteErr := svc.repo.Unfollow(ctx, user.ID, userToUnfollowID)
	if deleteErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			unfollowDefaultErr,
			fmt.Sprintf("Could not delete follow between user %d and user %d. Err: %v", user.ID, userToUnfollowID, deleteErr),
		)
	}

	return nil
}
