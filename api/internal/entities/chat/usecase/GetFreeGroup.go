package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetFreeGroup(ctx context.Context, channelID string) (*response.GetFreeGroup, error) {
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
			"user was nil in GetFreeGroup",
		)
	}

	runnable := svc.repo.MasterNode()
	group, getGroupErr := svc.repo.GetFreeGroup(ctx, runnable, channelID)
	if getGroupErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong getting group details.",
			fmt.Sprintf("Error getting group %s: %v", channelID, getGroupErr),
		)
	} else if group == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusBadRequest,
			"No group found.",
			"No free group found.",
		)
	}

	createdByUser, getCreatedByUserErr := svc.user.GetUserByID(ctx, group.CreatedByUserID)
	if getCreatedByUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error getting user %d: %v", group.CreatedByUserID, getCreatedByUserErr),
		)
	}
	if createdByUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"createdByUser was nil",
		)
	}

	fillGroupResponseErr := svc.fillFreeGroupResponse(ctx, user.ID, group)
	if fillGroupResponseErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong getting group details.",
			fmt.Sprintf("Error filling group response %s: %v", channelID, fillGroupResponseErr),
		)
	}

	res := &response.GetFreeGroup{
		CreatedByUser: createdByUser,
		Group:         group,
	}
	return res, nil
}

func (svc *usecase) fillFreeGroupResponse(ctx context.Context, userID int, group *response.FreeGroup) error {
	getGroupChannelParams := sendbird.GetGroupChannelParams{
		ShowMember: true,
	}
	channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(group.ChannelID, getGroupChannelParams)
	if getChannelErr != nil {
		vlog.Errorf(ctx, "Error getting sendbird channel %s in fillFreeGroupResponse: %v", group.ChannelID, getChannelErr)
		return getChannelErr
	}

	group.Channel = channel
	isMemberResponse, isMemberErr := svc.sendbirdClient.IsMemberInGroupChannel(channel.ChannelURL, fmt.Sprint(userID))
	if isMemberErr != nil {
		vlog.Errorf(ctx, "Error checking is member for sendbird channel %s in fillFreeGroupResponse: %v", group.ChannelID, isMemberErr)
		return isMemberErr
	}
	group.IsMember = isMemberResponse.IsMember
	return nil
}
