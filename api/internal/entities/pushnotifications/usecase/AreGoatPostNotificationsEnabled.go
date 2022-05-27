package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) AreGoatPostNotificationsEnabled(ctx context.Context, goatUserID int) (bool, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return false, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if user == nil {
		return false, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"user was nil in AreGoatPostNotificationsEnabled",
		)
	}

	enabled, areNotificationsEnabledErr := svc.repo.AreGoatPostNotificationsEnabled(ctx, user.ID, goatUserID)
	if areNotificationsEnabledErr != nil {
		return false, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error checking AreGoatPostNotificationsEnabled: %v", areNotificationsEnabledErr),
		)
	}

	return enabled, nil
}
