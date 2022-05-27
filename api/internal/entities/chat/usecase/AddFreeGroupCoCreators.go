package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errAddingFreeGroupCoCreators = "Something went wrong adding other creators."

func (svc *usecase) AddFreeGroupCoCreators(ctx context.Context, req request.AddFreeGroupChatCoCreators) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errAddingFreeGroupCoCreators,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errAddingFreeGroupCoCreators,
			"user was nil in AddFreeGroupCoCreators",
		)
	}

	if user.Type != "GOAT" {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can have co-creators.",
			"Only creators can have co-creators.",
		)
	}

	groupInfo, getGroupInfoErr := svc.GetFreeGroup(ctx, req.ChannelID)
	if getGroupInfoErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errAddingFreeGroupCoCreators,
			fmt.Sprintf("Error getting free group %s: %v", req.ChannelID, getGroupInfoErr),
		)
	} else if groupInfo == nil {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"No group found.",
			fmt.Sprintf("No free group found for channel ID %s", req.ChannelID),
		)
	}

	if groupInfo.CreatedByUser.ID != user.ID {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only the owner of this group can add co-creators.",
			"Only the owner of this group can add co-creators.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errAddingFreeGroupCoCreators,
			fmt.Sprintf("Error creating transaction: %v", txErr),
		)
	}
	defer tx.Rollback(ctx)

	addCoCreatorsErr := svc.repo.AddFreeGroupChatCoCreators(ctx, tx, req)
	if addCoCreatorsErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errAddingFreeGroupCoCreators,
			fmt.Sprintf("Error adding cocreators: %v", addCoCreatorsErr),
		)
	}

	operatorIDs := []string{}
	for _, sendbirdUser := range groupInfo.Group.Channel.Operators {
		operatorIDs = append(operatorIDs, sendbirdUser.UserID)
	}
	for _, userID := range req.CoCreators {
		operatorIDs = append(operatorIDs, fmt.Sprint(userID))
	}

	updateGroupReq := &sendbird.UpdateGroupChannelParams{
		OperatorIDs: operatorIDs,
	}
	upsertData := false // no Data field changes
	_, updateChannelErr := svc.sendbirdClient.UpdateGroupChannel(req.ChannelID, updateGroupReq, upsertData)
	if updateChannelErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errAddingFreeGroupCoCreators,
			fmt.Sprintf("Error updating sendbird channel %s: %v", req.ChannelID, updateChannelErr),
		)
	}

	inviteMembersReq := &sendbird.InviteGroupChannelMembers{
		UserIDs: operatorIDs,
	}
	_, inviteMembersErr := svc.sendbirdClient.InviteGroupChannelMembers(req.ChannelID, inviteMembersReq)
	if inviteMembersErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errAddingFreeGroupCoCreators,
			fmt.Sprintf("Error inviting co-creators to sendbird channel %s: %v", req.ChannelID, inviteMembersErr),
		)
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errAddingFreeGroupCoCreators,
			fmt.Sprintf("Error committing transaction: %v", commitErr),
		)
	}
	return nil
}
