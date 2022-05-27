package service

import (
	"context"
	"fmt"
	"net/http"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const batchUnsubscribeErr = "Something went wrong when trying to unsubscribe in BatchUnsubscribePaidGroup."

func (svc *usecase) BatchUnsubscribePaidGroup(ctx context.Context, req cloudtasks.StripePaidGroupUnsubscribeTask) error {
	subscriptions := []response.PaidGroupChatSubscription{}
	for _, userID := range req.UserIDs {
		// The creator shouldn't be subscribed to themselves
		if req.GoatUserID == userID {
			continue
		}

		subscription, getSubErr := svc.repo.GetPaidGroupSubscription(ctx, svc.repo.MasterNode(), userID, req.ChannelID)
		if getSubErr != nil {
			return httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				batchUnsubscribeErr,
				fmt.Sprintf("Error getting paid group chat subscription for user %d and channel %s: %v", userID, req.ChannelID, getSubErr),
			)
		} else if subscription == nil {
			vlog.Errorf(ctx, "Paid group subscription for user %d and channel %s does not exist. This may be a duplicate task.", userID, req.ChannelID)
		} else {
			subscriptions = append(subscriptions, *subscription)
		}
	}

	var unsubFinalErr error
	for _, oldSubscription := range subscriptions {
		unsubErr := svc.unsubscribe(ctx, &oldSubscription)
		if unsubErr != nil {
			unsubFinalErr = unsubErr
		}
	}

	if unsubFinalErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			batchUnsubscribeErr,
			fmt.Sprintf("Error unsubscribing from paid group chat in BatchUnsubscribePaidGroup for channel %s: %v", req.ChannelID, unsubFinalErr),
		)
	}

	return unsubFinalErr
}

func (svc *usecase) unsubscribe(ctx context.Context, oldSubscription *response.PaidGroupChatSubscription) error {
	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return txErr
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

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
		return updateSubErr
	}

	deleteSubErr := svc.repo.DeletePaidGroupSubscription(ctx, tx, oldSubscription.ChannelID, oldSubscription.UserID)
	if deleteSubErr != nil {
		return deleteSubErr
	}

	_, cancelSubErr := svc.stripeClient.Subscriptions.Cancel(oldSubscription.StripeSubscriptionID, nil)
	if cancelSubErr != nil {
		return cancelSubErr
	}

	commit = true
	return nil
}
