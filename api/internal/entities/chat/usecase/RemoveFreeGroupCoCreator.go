package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errRemovingFreeGroupCoCreator = "Something went wrong removing other creators."

func (svc *usecase) RemoveFreeGroupCoCreator(ctx context.Context, req request.RemoveFreeGroupChatCoCreator) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errRemovingFreeGroupCoCreator,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errRemovingFreeGroupCoCreator,
			"user was nil in RemoveFreeGroupCoCreator",
		)
	}

	if user.Type != "GOAT" {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can remove co-creators.",
			"Only creators can remove co-creators.",
		)
	}

	groupInfo, getGroupInfoErr := svc.GetFreeGroup(ctx, req.ChannelID)
	if getGroupInfoErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errRemovingFreeGroupCoCreator,
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
			"Only the owner of this group can remove co-creators.",
			"Only the owner of this group can remove co-creators.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errRemovingFreeGroupCoCreator,
			fmt.Sprintf("Error creating transaction: %v", txErr),
		)
	}
	defer tx.Rollback(ctx)

	addCoCreatorErr := svc.repo.RemoveFreeGroupChatCoCreator(ctx, tx, req)
	if addCoCreatorErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errRemovingFreeGroupCoCreator,
			fmt.Sprintf("Error removing cocreator: %v", addCoCreatorErr),
		)
	}

	operatorIDs := []string{}
	deleteUserIDStr := fmt.Sprint(req.UserID)
	for _, sendbirdUser := range groupInfo.Group.Channel.Operators {
		if sendbirdUser.UserID != deleteUserIDStr {
			operatorIDs = append(operatorIDs, sendbirdUser.UserID)
		}
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
			errRemovingFreeGroupCoCreator,
			fmt.Sprintf("Error updating sendbird channel %s: %v", req.ChannelID, updateChannelErr),
		)
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errRemovingFreeGroupCoCreator,
			fmt.Sprintf("Error committing transaction: %v", commitErr),
		)
	}
	return nil
}
