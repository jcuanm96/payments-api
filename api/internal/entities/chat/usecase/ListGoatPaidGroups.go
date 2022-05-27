package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errListingCreatorPaidGroups = "Something went wrong trying to list creator's paid groups."

func (svc *usecase) ListGoatPaidGroups(ctx context.Context, req request.ListGoatPaidGroups) (*response.ListGoatPaidGroups, error) {
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
			"user was nil in ListGoatPaidGroups",
		)
	}

	goatUser, getGoatUserErr := svc.user.GetUserByID(ctx, req.GoatID)
	if getGoatUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errListingCreatorPaidGroups,
			fmt.Sprintf("Error getting creator user %d: %v", req.GoatID, getGoatUserErr),
		)
	} else if goatUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errListingCreatorPaidGroups,
			fmt.Sprintf("Creator user %d was nil", req.GoatID),
		)
	} else if goatUser.Type != "GOAT" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can have paid groups.",
			"Only creators can have paid groups.",
		)
	}

	groups, listGroupsErr := svc.repo.ListGoatPaidGroups(ctx, req.GoatID, req.CursorID, req.Limit)
	if listGroupsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errListingCreatorPaidGroups,
			fmt.Sprintf("Error listing paid groups for creator %d: %v", req.GoatID, listGroupsErr),
		)
	}

	var wg sync.WaitGroup
	wg.Add(len(groups))
	for i := range groups {
		go func(group *response.PaidGroup) {
			defer wg.Done()
			svc.fillPaidGroupResponse(ctx, user.ID, group)
		}(&groups[i])
	}

	wg.Wait()

	res := &response.ListGoatPaidGroups{
		GoatUser: *goatUser,
		Groups:   groups,
	}

	return res, nil
}

func (svc *usecase) fillPaidGroupResponse(ctx context.Context, userID int, group *response.PaidGroup) error {
	getGroupChannelParams := sendbird.GetGroupChannelParams{
		ShowMember: true,
	}
	channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(group.ChannelID, getGroupChannelParams)
	if getChannelErr != nil {
		vlog.Errorf(ctx, "Error getting sendbird channel %s in fillPaidGroupResponse: %v", group.ChannelID, getChannelErr)
		return getChannelErr
	}

	group.Channel = channel
	isMemberResponse, isMemberErr := svc.sendbirdClient.IsMemberInGroupChannel(channel.ChannelURL, fmt.Sprint(userID))
	if isMemberErr != nil {
		vlog.Errorf(ctx, "Error checking is member for sendbird channel %s in fillPaidGroupResponse: %v", group.ChannelID, isMemberErr)
		return isMemberErr
	}
	group.IsMember = isMemberResponse.IsMember
	return nil
}
