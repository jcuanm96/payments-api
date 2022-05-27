package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errUpdatingPushSetting = "Something went wrong updating push notification setting."

func (svc *usecase) UpdateSetting(ctx context.Context, req request.UpdatePushSetting) error {
	user, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPushSetting,
			fmt.Sprintf("Could not get current user: %v", getUserErr),
		)
	} else if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errUpdatingPushSetting,
			"current user was nil in UpdateSettings",
		)
	}

	updateSettingsErr := svc.repo.UpdateSetting(ctx, svc.repo.MasterNode(), user.ID, req)
	if updateSettingsErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpdatingPushSetting,
			fmt.Sprintf("Error updating push settings: %v", updateSettingsErr),
		)
	}

	return nil
}
