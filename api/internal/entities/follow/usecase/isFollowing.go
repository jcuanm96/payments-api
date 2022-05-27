package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) IsFollowing(ctx context.Context, goatUserID int) (*response.IsFollowing, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)

	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"Current user was nil in IsFollowing",
		)
	}

	isFollowing, isFollowingErr := svc.repo.IsFollowing(ctx, user.ID, goatUserID)
	if isFollowingErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error checking IsFollowing: %v", isFollowingErr),
		)
	}

	return &response.IsFollowing{IsFollowing: isFollowing}, nil
}
