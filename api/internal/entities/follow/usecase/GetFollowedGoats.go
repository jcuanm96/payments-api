package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getFollowedCreatorsDefaultErr = "Something went wrong when trying to get the creators you're following."

func (svc *usecase) GetFollowedGoats(ctx context.Context, req request.GetFollowedGoats) (*response.GetFollowedGoats, error) {
	user, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getFollowedCreatorsDefaultErr,
			fmt.Sprintf("Could not find user in the current context. Err: %v", getUserErr),
		)
	} else if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			getFollowedCreatorsDefaultErr,
			"user returned nil for GetFollowedGoats.",
		)
	}

	goats, getFollowedCreatorsErr := svc.repo.GetFollowedGoats(ctx, user.ID, req)
	if getFollowedCreatorsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getFollowedCreatorsDefaultErr,
			fmt.Sprintf("Could not list creators for user %d. Err: %v", user.ID, getFollowedCreatorsErr),
		)
	}

	res := &response.GetFollowedGoats{
		Goats: goats,
	}

	return res, nil
}
