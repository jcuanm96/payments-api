package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errLeavingFreeGroup = "Something went wrong trying to leave group."

func (svc *usecase) LeaveFreeGroup(ctx context.Context, req request.LeaveGroup) error {
	currUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errLeavingFreeGroup,
			fmt.Sprintf("Could not find user in the current context: %v", getCurrUserErr),
		)
	} else if currUser == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errLeavingFreeGroup,
			"user was nil in LeaveFreeGroup",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errLeavingFreeGroup,
			fmt.Sprintf("Something went wrong when creating transaction: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	group, getGroupErr := svc.repo.GetFreeGroup(ctx, tx, req.ChannelID)
	if getGroupErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errLeavingFreeGroup,
			fmt.Sprintf("Error getting free group %s: %v", req.ChannelID, getGroupErr),
		)
	} else if group == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"No group found.",
			"No free group found.",
		)
	}

	if currUser.ID == group.CreatedByUserID {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You cannot leave your own group.",
			fmt.Sprintf("user %d tried to leave their own group %s", currUser.ID, req.ChannelID),
		)
	}

	leaveGroupChatParams := sendbird.LeaveGroupChannelParams{
		UserIDs: []string{fmt.Sprint(currUser.ID)},
	}

	leaveGroupChatErr := svc.sendbirdClient.LeaveGroupChannel(req.ChannelID, &leaveGroupChatParams)
	if leaveGroupChatErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errLeavingFreeGroup,
			fmt.Sprintf("Error leaving Sendbird paid group chat for user %d and free group channel %s: %v", currUser.ID, req.ChannelID, leaveGroupChatErr),
		)
	}

	commit = true

	return nil
}
