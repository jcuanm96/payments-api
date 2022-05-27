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

func (svc *usecase) UpsertBio(ctx context.Context, req request.UpsertBio) (*response.UpsertBio, error) {
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
			"user was nil in UpdateBio",
		)
	} else if user.Type != "GOAT" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can have a bio.",
			"Only creators can have a bio.",
		)
	}

	upsertErr := svc.repo.UpsertBio(ctx, user.ID, req.TextContent)
	if upsertErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error setting bio for creator: %d: %v", user.ID, upsertErr),
		)
	}

	res := response.UpsertBio{}
	return &res, nil
}
