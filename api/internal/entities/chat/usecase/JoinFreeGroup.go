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
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errJoiningFreeGroup = "Something went wrong when trying to join group. Please try again."

func (svc *usecase) JoinFreeGroup(ctx context.Context, req request.JoinFreeGroup) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"user was nil in JoinFreeGroup",
		)
	}

	var wg sync.WaitGroup

	var groupDBInfo *response.FreeGroup
	var getGroupErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		groupDBInfo, getGroupErr = svc.repo.GetFreeGroup(ctx, svc.repo.MasterNode(), req.ChannelID)
	}()

	var channel *sendbird.GroupChannel
	var getChannelErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		channel, getChannelErr = svc.sendbirdClient.GetGroupChannel(req.ChannelID, sendbird.GetGroupChannelParams{})
	}()

	var bannedUserID *int
	var getBannedUserErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		bannedUserID, getBannedUserErr = svc.repo.GetBannedChatUser(ctx, svc.repo.MasterNode(), user.ID, req.ChannelID)
	}()

	wg.Wait()

	if getGroupErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errJoiningFreeGroup,
			fmt.Sprintf("Error getting free group %s when trying to join: %v", req.ChannelID, getGroupErr),
		)
	} else if groupDBInfo == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"No group found.",
			"No group found.",
		)
	}

	if getChannelErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errJoiningFreeGroup,
			fmt.Sprintf("Error getting sendbird channel %s when trying to join: %v", req.ChannelID, getChannelErr),
		)
	}

	if getBannedUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errJoiningFreeGroup,
			fmt.Sprintf("Error checking if user %d is banned from free group %s: %v", user.ID, req.ChannelID, getBannedUserErr),
		)
	} else if bannedUserID != nil {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You are unable to join this group chat.",
			"You are unable to join this free group chat.",
		)
	}

	if groupDBInfo.IsMemberLimitEnabled && groupDBInfo.MemberLimit <= channel.MemberCount {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"This group chat is full.",
			"This free group chat is full.",
		)
	}

	joinParams := &sendbird.JoinGroupChannelParams{UserID: fmt.Sprint(user.ID)}
	joinGroupChannelErr := svc.sendbirdClient.JoinGroupChannel(req.ChannelID, joinParams)
	if joinGroupChannelErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errJoiningFreeGroup,
			fmt.Sprintf("Error user %d joining free group channel: %v", user.ID, joinGroupChannelErr),
		)
	}

	return nil
}
