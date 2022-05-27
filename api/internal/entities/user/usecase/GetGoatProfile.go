package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetGoatProfile(ctx context.Context, req request.GetGoatProfile) (*response.GetGoatProfile, error) {
	user, userErr := svc.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"Current user was nil in GetGoatProfile",
		)
	}

	res, getGoatProfileErr := svc.repo.GetGoatProfile(ctx, req.GoatID, user.ID)
	if getGoatProfileErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong getting creator's profile. Please try again.",
			fmt.Sprintf("Error getting creator profile: %v", getGoatProfileErr),
		)
	}

	return res, nil
}
