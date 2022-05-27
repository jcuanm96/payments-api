package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription/repositories"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errUpdatingPaidGroupInfo = "Something went wrong when trying to update paid group information."

func (svc *usecase) UpdatePaidGroupMetadata(ctx context.Context, req request.UpdateGroupMetadata) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroupInfo,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errUpdatingPaidGroupInfo,
			fmt.Sprintf("UpdatePaidGroupMetadata returned nil user for request: %v", req),
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
			errUpdatingPaidGroupInfo,
			fmt.Sprintf("Error getting existing paid group info for channel %s: %v", req.ChannelID, getPaidGroupInfoErr),
		)
	}

	if user.ID != oldPaidGroupInfo.GoatID {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You can't update this paid group chat's information.",
			"You can't update this paid group chat's information.",
		)
	}

	metadata := response.GroupMetadata{
		Description: req.Description,
		Benefit1:    req.Benefit1,
		Benefit2:    req.Benefit2,
		Benefit3:    req.Benefit3,
	}

	metadataBytes, marshalErr := json.Marshal(metadata)
	if marshalErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroupInfo,
			fmt.Sprintf("Error marshalling paid group chat info during update information: %v", marshalErr),
		)
	}

	_, updatePaidGroupErr := svc.repo.UpsertPaidGroupChat(
		ctx,
		user.ID,
		req.ChannelID,
		oldPaidGroupInfo.StripeProductID,
		int(oldPaidGroupInfo.PriceInSmallestDenom),
		oldPaidGroupInfo.Currency,
		oldPaidGroupInfo.LinkSuffix,
		oldPaidGroupInfo.MemberLimit,
		oldPaidGroupInfo.IsMemberLimitEnabled,
		metadataBytes,
		svc.repo.MasterNode(),
	)
	if updatePaidGroupErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPaidGroupInfo,
			fmt.Sprintf("Error updating paid group chat information: %v", updatePaidGroupErr),
		)
	}

	return nil
}
