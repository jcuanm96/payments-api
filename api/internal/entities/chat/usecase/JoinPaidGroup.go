package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errJoiningPaidGroup = "Something went wrong when trying to join paid group.  Please try again."

func (svc *usecase) JoinPaidGroup(ctx context.Context, req request.JoinPaidGroup) (*response.PaidGroupChatSubscription, error) {
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
			"user was nil in JoinPaidGroup",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errJoiningPaidGroup,
			fmt.Sprintf("Error starting transaction: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	var wg sync.WaitGroup

	var groupDBInfo *response.PaidGroup
	var getGroupErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		groupDBInfo, getGroupErr = svc.repo.GetPaidGroup(ctx, tx, req.ChannelID)
	}()

	var channel *sendbird.GroupChannel
	var getChannelErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		getGroupChannelParams := sendbird.GetGroupChannelParams{
			ShowMember: true,
		}
		channel, getChannelErr = svc.sendbirdClient.GetGroupChannel(req.ChannelID, getGroupChannelParams)
	}()

	wg.Wait()

	if getGroupErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong joining paid group.",
			fmt.Sprintf("Error getting paid group %s when trying to join: %v", req.ChannelID, getGroupErr),
		)
	} else if groupDBInfo == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"No paid group found.",
			"No paid group found.",
		)
	}

	if getChannelErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong joining paid group.",
			fmt.Sprintf("Error getting sendbird channel %s when trying to join: %v", req.ChannelID, getChannelErr),
		)
	}

	bannedUserID, getBannedUserErr := svc.repo.GetBannedChatUser(ctx, tx, user.ID, req.ChannelID)
	if getBannedUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errJoiningPaidGroup,
			fmt.Sprintf("Error checking if user %d is banned from paid group %s: %v", user.ID, req.ChannelID, getBannedUserErr),
		)
	} else if bannedUserID != nil {
		return nil, httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You are unable to join this paid group chat.",
			"You are unable to join this paid group chat.",
		)
	}

	if groupDBInfo.IsMemberLimitEnabled && groupDBInfo.MemberLimit <= channel.MemberCount {
		return nil, httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"This paid group chat is full.",
			"This paid group chat is full.",
		)
	}

	newSubscription, subscribeHttpErr := svc.subscription.SubscribePaidGroup(ctx, tx, req.ChannelID, user)
	if subscribeHttpErr != nil {
		vlog.Errorf(ctx, "Error subscribing user %d to paid group %s: %v", user.ID, req.ChannelID, subscribeHttpErr)
		return nil, subscribeHttpErr
	}

	joinParams := &sendbird.JoinGroupChannelParams{UserID: fmt.Sprint(user.ID)}
	joinGroupChannelErr := svc.sendbirdClient.JoinGroupChannel(req.ChannelID, joinParams)
	if joinGroupChannelErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errJoiningPaidGroup,
			fmt.Sprintf("Error user %d joining paid group channel: %v", user.ID, joinGroupChannelErr),
		)
	}

	commit = true
	return newSubscription, nil
}
