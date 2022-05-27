package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errUpdatingFreeGroupLink = "Something went wrong when trying to update your group chat link."

func (svc *usecase) UpdateFreeGroupLink(ctx context.Context, req request.UpdateGroupLink) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupLink,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errUpdatingFreeGroupLink,
			fmt.Sprintf("UpdateFreeGroupLink returned nil user for request: %v", req),
		)
	}

	oldFreeGroupInfo, getFreeGroupInfoErr := svc.repo.GetFreeGroup(ctx, svc.repo.MasterNode(), req.ChannelID)
	if getFreeGroupInfoErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupLink,
			fmt.Sprintf("Error getting existing free group info for channel %s: %v", req.ChannelID, getFreeGroupInfoErr),
		)
	}

	if user.ID != oldFreeGroupInfo.CreatedByUserID {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You can't update this group chat's link.",
			"You can't update this free group chat's link.",
		)
	}

	metadataBytes, marshalErr := json.Marshal(oldFreeGroupInfo.Metadata)
	if marshalErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupLink,
			fmt.Sprintf("Error marshalling free group chat metadata during update link: %v", marshalErr),
		)
	}

	isTakenReq := request.CheckLinkSuffixIsTaken{
		LinkSuffix: req.LinkSuffix,
		ChannelID:  &req.ChannelID,
	}

	takenRes, isTakenErr := svc.CheckLinkSuffixIsTaken(ctx, isTakenReq)
	if isTakenErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupLink,
			fmt.Sprintf("Error checking if link is taken: %v", isTakenErr),
		)
	}

	if takenRes.Taken {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"This link is unavailable.",
			"This link is unavailable.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupLink,
			fmt.Sprintf("Something went wrong when creating transaction. Err: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	_, updateLinkErr := svc.repo.UpsertFreeGroupChat(
		ctx,
		user.ID,
		req.ChannelID,
		req.LinkSuffix,
		oldFreeGroupInfo.MemberLimit,
		oldFreeGroupInfo.IsMemberLimitEnabled,
		metadataBytes,
		tx,
	)
	if updateLinkErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupLink,
			fmt.Sprintf("Error updating free group chat link: %v", updateLinkErr),
		)
	}

	channelData := response.GroupChannelData{
		LinkSuffix: &req.LinkSuffix,
	}

	channelDatab, marshalErr := json.Marshal(channelData)
	if marshalErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupLink,
			fmt.Sprintf("Error marshaling channel data in UpdateFreeGroupLink: %v", marshalErr),
		)
	}

	updateChannelParams := sendbird.UpdateGroupChannelParams{
		Data: string(channelDatab),
	}
	shouldUpsertData := true
	_, updateChannelDataErr := svc.sendbirdClient.UpdateGroupChannel(req.ChannelID, &updateChannelParams, shouldUpsertData)
	if updateChannelDataErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupLink,
			fmt.Sprintf("Error updating channel data UpdateFreeGroupLink: %v", updateChannelDataErr),
		)
	}

	commit = true
	return nil
}
