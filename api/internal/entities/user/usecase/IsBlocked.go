package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) IsBlocked(ctx context.Context, userID int) (bool, error) {
	currUser, userErr := svc.GetCurrentUser(ctx)
	if userErr != nil {
		return false, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if currUser == nil {
		return false, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"Current user was nil in IsBlocked",
		)
	}

	isBlocked, isBlockedErr := svc.repo.IsBlocked(ctx, currUser.ID, userID)
	if isBlockedErr != nil {
		return false, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error checking user %d is blocked by %d: %v", userID, currUser.ID, isBlockedErr),
		)
	}

	return isBlocked, nil
}
