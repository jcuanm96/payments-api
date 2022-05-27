package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/webhooks"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

func (svc *usecase) StripeSubscriptionPaymentFailure(ctx context.Context, event stripe.Event) (*response.HandleSubscriptionPaymentEvent, error) {
	subscriptionID := event.Data.Object["subscription"].(string)
	getSubscriptionParams := stripe.SubscriptionParams{}
	subscription, getSubscriptionErr := svc.stripeClient.Subscriptions.Get(subscriptionID, &getSubscriptionParams)

	if getSubscriptionErr != nil {
		vlog.Errorf(ctx, "Error getting subscription from Stripe for subscription %s. Err: %s", subscription.ID, getSubscriptionErr.Error())
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to get subscription from Stripe event",
		)
	}

	metadata, subscriptionMetadataErr := webhooks.StripeValidateSubscriptionMetadata(ctx, subscription.Metadata, subscription.ID)
	if subscriptionMetadataErr != nil {
		vlog.Errorf(ctx, "Error validating subscription %s metadata %s", subscription.ID, subscriptionMetadataErr.Error())
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when validating subscription metadata",
		)
	}

	// If the metadata contains a channelID, we know it's a paid group chat.
	if metadata.ChannelID != nil {
		newSubscription := &response.PaidGroupChatSubscription{
			StripeSubscriptionID: "",
			IsRenewing:           false,
			CurrentPeriodEnd:     time.Now(),
			UserID:               metadata.CustomerUserID,
			GoatUser:             response.User{ID: metadata.ProviderUserID},
			ChannelID:            *metadata.ChannelID,
		}

		upsertErr := svc.subscription.UpsertPaidGroupSubscription(ctx, svc.repo.MasterNode(), newSubscription)
		if upsertErr != nil {
			vlog.Errorf(ctx, "Error upserting paid group subscription %s into db: %s", subscription.ID, upsertErr.Error())
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when trying to upsert paid group subscription.",
			)
		}

		leaveGroupParams := &sendbird.LeaveGroupChannelParams{
			UserIDs: []string{fmt.Sprint(metadata.CustomerUserID)},
		}
		leaveGroupErr := svc.sendbirdClient.LeaveGroupChannel(*metadata.ChannelID, leaveGroupParams)
		if leaveGroupErr != nil {
			vlog.Errorf(ctx, "Error leaving group channel %s: %s", *metadata.ChannelID, leaveGroupErr.Error())
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when trying to leave group channel.",
			)
		}
	} else if metadata.TierID != nil { // If the metadata contains a tierID, we know it's a user subscription.
		newSubscription := response.UserSubscription{
			StripeSubscriptionID: "",
			IsRenewing:           false,
			CurrentPeriodEnd:     time.Now(),
			UserID:               metadata.CustomerUserID,
			GoatUser:             response.User{ID: metadata.ProviderUserID},
			TierID:               *metadata.TierID,
		}

		upsertErr := svc.subscription.UpsertUserSubscription(ctx, svc.repo.MasterNode(), newSubscription)
		if upsertErr != nil {
			vlog.Errorf(ctx, "Error upserting user subscription %s into db: %s", subscription.ID, upsertErr.Error())
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when trying to upsert user's subscription",
			)
		}
	}

	_, cancelErr := svc.stripeClient.Subscriptions.Cancel(subscriptionID, nil)
	if cancelErr != nil {
		vlog.Errorf(ctx, "Error cancelling subscription in Stripe for subscription %s: %s", subscriptionID, cancelErr.Error())
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when cancelling subscription",
		)
	}

	res := response.HandleSubscriptionPaymentEvent{}

	return &res, nil
}
