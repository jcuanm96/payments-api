package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) UpsertUserSubscription(ctx context.Context, runnable utils.Runnable, subscription response.UserSubscription) error {
	updateSubscriptionErr := svc.repo.UpsertUserSubscription(ctx, runnable, subscription)
	if updateSubscriptionErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong updating subscription status",
			fmt.Sprintf("Error updating subscription status for user %d and creator %d: %v", subscription.UserID, subscription.GoatUser.ID, updateSubscriptionErr),
		)
	}
	return nil
}
