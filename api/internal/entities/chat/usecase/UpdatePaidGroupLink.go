package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errUpdatingPaidGroup = "Something went wrong when trying to update your paid group chat"

func (svc *usecase) UpdatePaidGroupLink(ctx context.Context, req request.UpdateGroupLink) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroup,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errUpdatingPaidGroup,
			fmt.Sprintf("UpdatePaidGroupLink returned nil user for request: %v", req),
		)
	} else if user.Type != "GOAT" {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators are allowed to have paid group chats.",
			"Only creators are allowed to have paid group chats.",
		)
	}

	oldPaidGroupInfo, getPaidGroupInfoErr := repositories.GetPaidGroupProductInfo(ctx, svc.repo.MasterNode(), req.ChannelID)
	if getPaidGroupInfoErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroup,
			fmt.Sprintf("Error getting existing paid group info for channel %s: %v", req.ChannelID, getPaidGroupInfoErr),
		)
	}

	if user.ID != oldPaidGroupInfo.GoatID {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You can't update this paid group chat's link.",
			"You can't update this paid group chat's link.",
		)
	}

	metadataBytes, marshalErr := json.Marshal(oldPaidGroupInfo.Metadata)
	if marshalErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroup,
			fmt.Sprintf("Error marshalling paid group chat info during update link: %v", marshalErr),
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
			errUpdatingPaidGroup,
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
			errUpdatingPaidGroup,
			fmt.Sprintf("Something went wrong when creating transaction. Err: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	_, updateLinkErr := svc.repo.UpsertPaidGroupChat(
		ctx,
		user.ID,
		req.ChannelID,
		oldPaidGroupInfo.StripeProductID,
		int(oldPaidGroupInfo.PriceInSmallestDenom),
		oldPaidGroupInfo.Currency,
		req.LinkSuffix,
		oldPaidGroupInfo.MemberLimit,
		oldPaidGroupInfo.IsMemberLimitEnabled,
		metadataBytes,
		tx,
	)
	if updateLinkErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroup,
			fmt.Sprintf("Error updating paid group chat link: %v", updateLinkErr),
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
			"Something went wrong when trying to update your paid group chat",
			fmt.Sprintf("Error marshaling channel data in UpdatePaidGroupLink: %v", marshalErr),
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
			"Something went wrong when trying to update your paid group chat",
			fmt.Sprintf("Error updating channel data UpdatePaidGroupLink: %v", updateChannelDataErr),
		)
	}

	commit = true

	return nil
}
