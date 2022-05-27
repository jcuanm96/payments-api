package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetPaidGroup(ctx context.Context, channelID string) (*response.GetPaidGroup, error) {
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
			"user was nil in GetPaidGroup",
		)
	}

	runnable := svc.repo.MasterNode()
	group, getGroupErr := svc.repo.GetPaidGroup(ctx, runnable, channelID)
	if getGroupErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong getting paid group details.",
			fmt.Sprintf("Error getting paid group %s: %v", channelID, getGroupErr),
		)
	} else if group == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusBadRequest,
			"No paid group found.",
			"No paid group found.",
		)
	}

	goatUser, getGoatUserErr := svc.user.GetUserByID(ctx, group.GoatID)
	if getGoatUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error getting creator user %d: %v", group.GoatID, getGoatUserErr),
		)
	}
	if goatUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"goatUser was nil",
		)
	}

	fillGroupResponseErr := svc.fillPaidGroupResponse(ctx, user.ID, group)
	if fillGroupResponseErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong getting paid group details.",
			fmt.Sprintf("Error filling paid group response %s: %v", channelID, fillGroupResponseErr),
		)
	}

	res := &response.GetPaidGroup{
		GoatUser: goatUser,
		Group:    group,
	}
	return res, nil
}
