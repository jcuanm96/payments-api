package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription/repositories"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errUpdatingPaidGroupPrice = "Something went wrong when trying to update paid group price."

func (svc *usecase) UpdatePaidGroupPrice(ctx context.Context, req request.UpdatePaidGroupPrice) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroupPrice,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errUpdatingPaidGroupPrice,
			fmt.Sprintf("UpdatePaidGroupPrice returned nil user for request: %v", req),
		)
	} else if user.Type != "GOAT" {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators are allowed to have paid group chats.",
		)
	}

	oldPaidGroupInfo, getPaidGroupInfoErr := repositories.GetPaidGroupProductInfo(ctx, svc.repo.MasterNode(), req.ChannelID)
	if getPaidGroupInfoErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroupPrice,
			fmt.Sprintf("Error getting existing paid group info for channel %s: %v", req.ChannelID, getPaidGroupInfoErr),
		)
	}

	if user.ID != oldPaidGroupInfo.GoatID {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You can't update this paid group chat's price.",
		)
	}

	metadataBytes, marshalErr := json.Marshal(oldPaidGroupInfo.Metadata)
	if marshalErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroupPrice,
			fmt.Sprintf("Error marshalling paid group chat info during update price: %v", marshalErr),
		)
	}

	_, updatePriceErr := svc.repo.UpsertPaidGroupChat(
		ctx,
		user.ID,
		req.ChannelID,
		oldPaidGroupInfo.StripeProductID,
		req.PriceInSmallestDenom,
		req.Currency,
		oldPaidGroupInfo.LinkSuffix,
		oldPaidGroupInfo.MemberLimit,
		oldPaidGroupInfo.IsMemberLimitEnabled,
		metadataBytes,
		svc.repo.MasterNode(),
	)
	if updatePriceErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroupPrice,
			fmt.Sprintf("Error updating paid group chat price: %v", updatePriceErr),
		)
	}

	return nil
}
