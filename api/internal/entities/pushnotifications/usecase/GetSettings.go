package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errGettingPushSettings = "Something went wrong getting push notification settings."

func (svc *usecase) GetSettings(ctx context.Context) (*response.GetPushSettings, error) {
	user, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errGettingPushSettings,
			fmt.Sprintf("Could not find user in the current context: %v", getUserErr),
		)
	} else if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errGettingPushSettings,
			"user was nil in GetSettings",
		)
	}

	res, getSettingsErr := svc.repo.GetSettings(ctx, svc.repo.MasterNode(), user.ID)
	if getSettingsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errGettingPushSettings,
			fmt.Sprintf("Error getting settings: %v", getSettingsErr),
		)
	}

	return res, nil
}
