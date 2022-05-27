package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const unsubscribeCurrUserFromGoatErr = "Something went wrong when trying to unsubscribe from creator."

func (svc *usecase) UnsubscribeCurrUserFromGoat(ctx context.Context, req request.UnsubscribeCurrUserFromGoat) (*response.UserSubscription, error) {
	currUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			unsubscribeCurrUserFromGoatErr,
			fmt.Sprintf("Could not find user in the current context: %v", getCurrUserErr),
		)
	} else if currUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			unsubscribeCurrUserFromGoatErr,
		)
	}

	oldSubscription, getSubErr := svc.repo.GetUserSubscriptionByGoatID(ctx, currUser.ID, req.GoatUserID)
	if getSubErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			unsubscribeCurrUserFromGoatErr,
			fmt.Sprintf("Error getting subscription for user %d and creator %d: %v", currUser.ID, req.GoatUserID, getSubErr),
		)
	} else if oldSubscription == nil || !oldSubscription.IsRenewing {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You are not currently subscribed to this creator.",
			"You are not currently subscribed to this creator.",
		)
	} else if oldSubscription != nil && oldSubscription.IsRenewing && oldSubscription.StripeSubscriptionID == "" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"There was a problem unsubscribing. Please wait a few moments then try again.",
			"The user likely subscribed then unsubscribed quickly, wait for webhook to update StripeSubscriptionID",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			unsubscribeCurrUserFromGoatErr,
			fmt.Sprintf("Something went wrong when creating transaction. Err: %v", txErr),
		)
	}
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	newSubscription := response.UserSubscription{
		ID:                   oldSubscription.ID,
		UserID:               oldSubscription.UserID,
		GoatUser:             oldSubscription.GoatUser,
		TierID:               oldSubscription.TierID,
		CurrentPeriodEnd:     oldSubscription.CurrentPeriodEnd,
		StripeSubscriptionID: "",
		IsRenewing:           false,
	}

	updateSubErr := svc.repo.UpsertUserSubscription(ctx, tx, newSubscription)
	if updateSubErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			unsubscribeCurrUserFromGoatErr,
			fmt.Sprintf("Error updating subscription status for user %d and creator %d: %v", currUser.ID, req.GoatUserID, updateSubErr),
		)
	}

	_, cancelSubErr := svc.stripeClient.Subscriptions.Cancel(oldSubscription.StripeSubscriptionID, nil)
	if cancelSubErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			unsubscribeCurrUserFromGoatErr,
			fmt.Sprintf("Error cancelling subscription %s in Stripe for user %d and creator %d: %v", oldSubscription.StripeSubscriptionID, currUser.ID, req.GoatUserID, cancelSubErr),
		)
	}

	// Make new subscription object so we don't return isRenewing as true
	// and remove the StripeSubscriptionID
	returnedSubscription := response.UserSubscription{
		ID:               oldSubscription.ID,
		UserID:           oldSubscription.UserID,
		GoatUser:         oldSubscription.GoatUser, // This user only has ID set.
		CurrentPeriodEnd: oldSubscription.CurrentPeriodEnd,
		TierID:           oldSubscription.TierID,
		IsRenewing:       false,
	}
	commit = true
	return &returnedSubscription, nil
}
