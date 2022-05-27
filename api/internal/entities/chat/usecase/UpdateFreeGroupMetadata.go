package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errUpdatingFreeGroupInfo = "Something went wrong when trying to update group information."

func (svc *usecase) UpdateFreeGroupMetadata(ctx context.Context, req request.UpdateGroupMetadata) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupInfo,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errUpdatingFreeGroupInfo,
			fmt.Sprintf("UpdateFreeGroupMetadata returned nil user for request: %v", req),
		)
	}

	oldFreeGroupInfo, getFreeGroupInfoErr := svc.repo.GetFreeGroup(ctx, svc.repo.MasterNode(), req.ChannelID)
	if getFreeGroupInfoErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupInfo,
			fmt.Sprintf("Error getting existing free group info for channel %s: %v", req.ChannelID, getFreeGroupInfoErr),
		)
	}

	if user.ID != oldFreeGroupInfo.CreatedByUserID {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You can't update this public group chat's information.",
			"You can't update this public free group chat's information.",
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
			errUpdatingFreeGroupInfo,
			fmt.Sprintf("Error marshalling free group chat metadata during update: %v", marshalErr),
		)
	}

	_, updateFreeGroupErr := svc.repo.UpsertFreeGroupChat(
		ctx,
		user.ID,
		req.ChannelID,
		oldFreeGroupInfo.LinkSuffix,
		oldFreeGroupInfo.MemberLimit,
		oldFreeGroupInfo.IsMemberLimitEnabled,
		metadataBytes,
		svc.repo.MasterNode(),
	)
	if updateFreeGroupErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingFreeGroupInfo,
			fmt.Sprintf("Error updating free group chat metadata: %v", updateFreeGroupErr),
		)
	}

	return nil
}
