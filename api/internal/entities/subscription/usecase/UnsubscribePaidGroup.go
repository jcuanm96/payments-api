package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/jackc/pgx/v4"
)

const unsubscribePaidGroupErr = "Something went wrong when trying to cancel paid group subscription."

func (svc *usecase) UnsubscribePaidGroup(ctx context.Context, tx pgx.Tx, channelID string, user *response.User) (*response.PaidGroupChatSubscription, error) {
	oldSubscription, getSubErr := svc.repo.GetPaidGroupSubscription(ctx, tx, user.ID, channelID)
	if getSubErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			unsubscribePaidGroupErr,
			fmt.Sprintf("Error getting paid group chat subscription for user %d and channel %s: %v", user.ID, channelID, getSubErr),
		)
		// No subscription, return nil error so we leave the group in the calling function
	} else if oldSubscription == nil || !oldSubscription.IsRenewing {
		return nil, nil
	} else if oldSubscription != nil && oldSubscription.IsRenewing && oldSubscription.StripeSubscriptionID == "" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"There was a problem cancelling paid group subscription. Please wait a few moments then try again.",
			"The user likely subscribed then unsubscribed quickly, wait for webhook to update StripeSubscriptionID",
		)
	}

	newSubscription := response.PaidGroupChatSubscription{
		ID:                   oldSubscription.ID,
		UserID:               oldSubscription.UserID,
		GoatUser:             oldSubscription.GoatUser,
		ChannelID:            oldSubscription.ChannelID,
		CurrentPeriodEnd:     oldSubscription.CurrentPeriodEnd,
		StripeSubscriptionID: "",
		IsRenewing:           false,
	}

	updateSubErr := svc.UpsertPaidGroupSubscription(ctx, tx, &newSubscription)
	if updateSubErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			unsubscribePaidGroupErr,
			fmt.Sprintf("Error updating paid group chat subscription for user %d and channel %s: %v", user.ID, channelID, updateSubErr),
		)
	}

	_, cancelSubErr := svc.stripeClient.Subscriptions.Cancel(oldSubscription.StripeSubscriptionID, nil)
	if cancelSubErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			unsubscribePaidGroupErr,
			fmt.Sprintf("Error cancelling subscription %s in Stripe for user %d and paid group channel %s: %v", oldSubscription.StripeSubscriptionID, user.ID, channelID, cancelSubErr),
		)
	}

	return &newSubscription, nil
}
