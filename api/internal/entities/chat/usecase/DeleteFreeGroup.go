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

const errDeletingFreeGroup = "Something went wrong trying to delete group."

func (svc *usecase) DeleteFreeGroup(ctx context.Context, req request.DeleteFreeGroup) error {
	currUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingFreeGroup,
			fmt.Sprintf("Error getting user from current context: %v", getCurrUserErr),
		)
	} else if currUser == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errDeletingFreeGroup,
			"Could not find user in the current context",
		)
	}

	channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(req.ChannelID, sendbird.GetGroupChannelParams{})
	if getChannelErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingFreeGroup,
			fmt.Sprintf("Error getting Sendbird channel for group chat %s and user %d: %v", req.ChannelID, currUser.ID, getChannelErr),
		)
	}

	currUserIDStr := fmt.Sprint(currUser.ID)
	if !channel.HasOperator(currUserIDStr) {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You cannot delete this group chat.",
			"You cannot delete this free group chat.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingFreeGroup,
			fmt.Sprintf("Error begining transaction in DeleteFreeGroup: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	deleteProductErr := svc.repo.DeleteFreeGroupProduct(ctx, tx, req.ChannelID)
	if deleteProductErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingFreeGroup,
			fmt.Sprintf("Error deleting free group chat product for channel %s and user %d: %v", req.ChannelID, currUser.ID, deleteProductErr),
		)
	}

	deleteChannelErr := svc.sendbirdClient.DeleteGroupChannel(req.ChannelID)
	if deleteChannelErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingFreeGroup,
			fmt.Sprintf("Error deleting Sendbird channel for group chat %s and user %d: %v", req.ChannelID, currUser.ID, deleteChannelErr),
		)
	}

	commit = true

	return nil
}
