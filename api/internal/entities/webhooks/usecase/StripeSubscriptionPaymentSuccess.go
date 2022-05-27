package service

import (
	"context"
	"net/http"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/internal/entities/webhooks"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

func (svc *usecase) StripeSubscriptionPaymentSuccess(ctx context.Context, event stripe.Event) (*wallet.StripeEventMetadata, error) {
	subscriptionID := event.Data.Object["subscription"].(string)
	getSubscriptionParams := stripe.SubscriptionParams{}
	getSubscriptionParams.AddExpand("latest_invoice.payment_intent")
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

	isSubscriptionSuccessful := false
	// Unsubscribe the user in Stripe if we fail to insert into the db
	defer func() {
		if !isSubscriptionSuccessful {
			_, cancelErr := svc.stripeClient.Subscriptions.Cancel(subscription.ID, nil)
			if cancelErr != nil {
				vlog.Errorf(ctx, "Error cancelling subscription in Stripe for subscription %s. Err: %s", subscription.ID, cancelErr.Error())
			}
			_, refundErr := svc.wallet.RefundSubscription(ctx, subscription.LatestInvoice.PaymentIntent.ID)
			if refundErr != nil {
				vlog.Errorf(ctx, "Error refunding subscription %s and payment intent %s. Err: %s", subscription.ID, subscription.LatestInvoice.PaymentIntent.ID, refundErr.Error())
			}
		}
	}()

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
	if subscription.LatestInvoice == nil {
		vlog.Errorf(ctx, "Error the Stripe invoice was nil for subscription %s provider %d and customer %d.", subscription.ID, metadata.ProviderUserID, metadata.CustomerUserID)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"LatestInvoice was nil",
		)
	}

	currentPeriodEnd := time.Unix(subscription.CurrentPeriodEnd, 0)
	isTrial := subscription.TrialEnd != 0 && time.Now().Unix() < subscription.TrialEnd
	metadata.IsTrial = &isTrial
	if isTrial {
		currentPeriodEnd = time.Unix(subscription.TrialEnd, 0)
	}

	// If the metadata contains a channelID, we know it's a paid group chat.
	if metadata.ChannelID != nil {
		newSubscription := &response.PaidGroupChatSubscription{
			StripeSubscriptionID: subscription.ID,
			CurrentPeriodEnd:     currentPeriodEnd,
			UserID:               metadata.CustomerUserID,
			GoatUser:             response.User{ID: metadata.ProviderUserID},
			ChannelID:            *metadata.ChannelID,
			IsRenewing:           true,
		}

		upsertErr := svc.subscription.UpsertPaidGroupSubscription(ctx, svc.repo.MasterNode(), newSubscription)
		if upsertErr != nil {
			vlog.Errorf(ctx, "Error upserting paid group subscription %s into db: %s", subscription.ID, upsertErr.Error())
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when trying to upsert paid group subscription",
			)
		}
	} else if metadata.TierID != nil { // If the metadata contains a tierID, we know it's a user subscription.
		newSubscription := response.UserSubscription{
			StripeSubscriptionID: subscription.ID,
			CurrentPeriodEnd:     currentPeriodEnd,
			UserID:               metadata.CustomerUserID,
			GoatUser:             response.User{ID: metadata.ProviderUserID},
			TierID:               *metadata.TierID,
			IsRenewing:           true,
		}

		upsertErr := svc.subscription.UpsertUserSubscription(ctx, svc.repo.MasterNode(), newSubscription)
		if upsertErr != nil {
			vlog.Errorf(ctx, "Error upserting user subscription %s into db. Err: %s", subscription.ID, upsertErr.Error())
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when trying to upsert user's subscription",
			)
		}
	}

	isSubscriptionSuccessful = true
	return metadata, nil
}
