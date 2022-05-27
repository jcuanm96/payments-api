package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const followDefaultErr = "Something went wrong when trying to follow."

func (svc *usecase) Follow(ctx context.Context, goatUserID int) (*response.User, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			followDefaultErr,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	} else if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			followDefaultErr,
			"The user returned nil for the Follow usecase.",
		)
	}

	userToFollow, getUserErr := svc.user.GetUserByID(ctx, goatUserID)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			followDefaultErr,
			fmt.Sprintf("There was an issue fetching user %d by ID.  Err: %v", goatUserID, getUserErr),
		)
	} else if userToFollow == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"The user you tried to follow does not exist.",
			"The user to follow returned nil.",
		)
	} else if userToFollow.Type != "GOAT" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"The user you tried to follow is not a creator.",
			"The user you tried to follow is not a creator.",
		)
	}

	followErr := svc.repo.Follow(ctx, user.ID, goatUserID)
	if followErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			followDefaultErr,
			fmt.Sprintf("There was an issue creating a follow between users %d and user %d.  Err: %v", user.ID, goatUserID, followErr),
		)
	}

	return userToFollow, nil
}
