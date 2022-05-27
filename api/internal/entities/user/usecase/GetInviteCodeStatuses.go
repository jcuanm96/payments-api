package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetInviteCodeStatuses(ctx context.Context) (*response.GetInviteCodeStatuses, error) {
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
			"user was nil in GetInviteCodeStatuses",
		)
	}

	invites, getInvitesErr := svc.repo.GetMyInvites(ctx, user.ID)
	if getInvitesErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong getting invite code statuses. Please try again.",
			fmt.Sprintf("Error getting invite code statuses: %v", getInvitesErr),
		)
	}

	res := response.GetInviteCodeStatuses{
		Invites: invites,
	}

	return &res, nil
}
