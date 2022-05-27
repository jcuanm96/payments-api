package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errListingUserFreeGroups = "Something went wrong trying to list user's free groups."

func (svc *usecase) ListUserFreeGroups(ctx context.Context, req request.ListUserFreeGroups) (*response.ListUserFreeGroups, error) {
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
			"user was nil in ListCreatorFreeGroups",
		)
	}

	createdByUser, getCreatedByUserErr := svc.user.GetUserByID(ctx, req.UserID)
	if getCreatedByUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errListingUserFreeGroups,
			fmt.Sprintf("Error getting user %d: %v", req.UserID, getCreatedByUserErr),
		)
	} else if createdByUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errListingUserFreeGroups,
			fmt.Sprintf("CreatedByUser %d was nil", req.UserID),
		)
	}

	groups, listGroupsErr := svc.repo.ListUserFreeGroups(ctx, req.UserID, req.CursorID, req.Limit)
	if listGroupsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errListingUserFreeGroups,
			fmt.Sprintf("Error listing free groups for user %d: %v", req.UserID, listGroupsErr),
		)
	}

	var wg sync.WaitGroup
	wg.Add(len(groups))
	for i := range groups {
		go func(group *response.FreeGroup) {
			defer wg.Done()
			svc.fillFreeGroupResponse(ctx, user.ID, group)
		}(&groups[i])
	}

	wg.Wait()

	res := &response.ListUserFreeGroups{
		User:   *createdByUser,
		Groups: groups,
	}

	return res, nil
}
