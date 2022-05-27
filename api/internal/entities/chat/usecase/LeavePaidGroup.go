package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errLeavingPaidGroup = "Something went wrong when trying to leave paid group."

func (svc *usecase) LeavePaidGroup(ctx context.Context, req request.LeaveGroup) (*response.LeavePaidGroup, error) {
	currUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errLeavingPaidGroup,
			fmt.Sprintf("Could not find user in the current context: %v", getCurrUserErr),
		)
	} else if currUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errLeavingPaidGroup,
			"user was nil in LeavePaidGroup",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errLeavingPaidGroup,
			fmt.Sprintf("Something went wrong when creating transaction: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	group, getGroupErr := svc.repo.GetPaidGroup(ctx, tx, req.ChannelID)
	if getGroupErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong getting paid group details.",
			fmt.Sprintf("Error getting paid group %s: %v", req.ChannelID, getGroupErr),
		)
	} else if group == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"No paid group found.",
			"No paid group found.",
		)
	}

	if currUser.ID == group.GoatID {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You cannot leave your own paid group.",
			fmt.Sprintf("creator %d tried to leave their own group %s", currUser.ID, req.ChannelID),
		)
	}

	cancelledSubscription, unsubscribeErr := svc.subscription.UnsubscribePaidGroup(ctx, tx, req.ChannelID, currUser)
	if unsubscribeErr != nil {
		return nil, unsubscribeErr
	}

	leaveGroupChatParams := sendbird.LeaveGroupChannelParams{
		UserIDs: []string{fmt.Sprint(currUser.ID)},
	}

	leaveGroupChatErr := svc.sendbirdClient.LeaveGroupChannel(req.ChannelID, &leaveGroupChatParams)
	if leaveGroupChatErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errLeavingPaidGroup,
			fmt.Sprintf("Error leaving Sendbird paid group chat for user %d and paid group channel %s: %v", currUser.ID, req.ChannelID, leaveGroupChatErr),
		)
	}

	commit = true

	return &response.LeavePaidGroup{Subscription: cancelledSubscription}, nil
}
