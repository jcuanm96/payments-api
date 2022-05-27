package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const updateLimitsErr = "Something went wrong when trying to update paid group member limits."

func (svc *usecase) UpdatePaidGroupMemberLimits(ctx context.Context, req request.UpdatePaidGroupMemberLimits) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			updateLimitsErr,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			updateLimitsErr,
			fmt.Sprintf("UpdatePaidGroupMemberLimits returned nil user for request: %v", req),
		)
	} else if user.Type != "GOAT" {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators are allowed to update paid group chats.",
			"Only creators are allowed to update paid group chats.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroup,
			fmt.Sprintf("Error creating transaction: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	var wg sync.WaitGroup

	var groupProductInfo *subscription.PaidGroupChatInfo
	var getGroupProductInfoErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		groupProductInfo, getGroupProductInfoErr = repositories.GetPaidGroupProductInfo(ctx, tx, req.ChannelID)
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

	if getGroupProductInfoErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			updateLimitsErr,
			fmt.Sprintf("Error getting existing paid group info for channel %s: %v", req.ChannelID, getGroupProductInfoErr),
		)
	} else if groupProductInfo == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"Could not find the paid group chat.",
			fmt.Sprintf("Error groupProductInfo returned nil for channel %s.", req.ChannelID),
		)
	}

	if getChannelErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			updateLimitsErr,
			fmt.Sprintf("Error getting sendbird channel %s when trying to update member limits: %v", req.ChannelID, getChannelErr),
		)
	} else if channel == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"No paid group found.",
			"No paid group found.",
		)
	}

	if req.IsMemberLimitEnabled == nil {
		req.IsMemberLimitEnabled = &groupProductInfo.IsMemberLimitEnabled
	}

	if req.MemberLimit == nil {
		req.MemberLimit = &groupProductInfo.MemberLimit
	}

	if user.ID != groupProductInfo.GoatID {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You can't update this paid group chat's member limits.",
			"You can't update this paid group chat's member limits.",
		)
	} else if req.MemberLimit != nil && channel.MemberCount > *req.MemberLimit {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You can't have a member limit that is less than the number of members already in the chat.",
			"You can't have a member limit that is less than the number of members already in the chat.",
		)
	}

	metadataBytes, marshalErr := json.Marshal(groupProductInfo.Metadata)
	if marshalErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to update paid group member limits.",
			fmt.Sprintf("Error marshalling paid group chat info during update member limits: %v", marshalErr),
		)
	}

	_, updatePriceErr := svc.repo.UpsertPaidGroupChat(
		ctx,
		user.ID,
		req.ChannelID,
		groupProductInfo.StripeProductID,
		int(groupProductInfo.PriceInSmallestDenom),
		groupProductInfo.Currency,
		groupProductInfo.LinkSuffix,
		*req.MemberLimit,
		*req.IsMemberLimitEnabled,
		metadataBytes,
		tx,
	)
	if updatePriceErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to update paid group member limits.",
			fmt.Sprintf("Error updating paid group chat member limits: %v", updatePriceErr),
		)
	}

	commit = true

	return nil
}
