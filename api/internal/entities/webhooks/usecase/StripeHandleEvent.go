package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

func (svc *usecase) StripeHandleEvent(ctx context.Context, stripeEvent stripe.Event) error {
	isValidSubscription := isValidSubscriptionEvent(stripeEvent)
	extractedStripeEvent := wallet.StripeTxnEvent{
		CreatedAt: stripeEvent.Created,
		ID:        stripeEvent.Data.Object["id"].(string),
	}

	if stripeEvent.Type == "invoice.paid" && isValidSubscription {
		metadata, successHandlerErr := svc.StripeSubscriptionPaymentSuccess(ctx, stripeEvent)
		if successHandlerErr != nil {
			vlog.Errorf(ctx, "Error handling successful subscription payment: %s", successHandlerErr.Error())
		}
		extractedStripeEvent.Metadata = metadata
	} else if stripeEvent.Type == "invoice.payment_failed" && isValidSubscription {
		_, failureHandlerErr := svc.StripeSubscriptionPaymentFailure(ctx, stripeEvent)
		if failureHandlerErr != nil {
			vlog.Errorf(ctx, "Error handling failed subscription payment: %v", failureHandlerErr)
		}
		return failureHandlerErr
	}

	if isValidPendingChargeEvent(stripeEvent) {
		return svc.StripePendingChargeEvent(ctx, stripeEvent)
	}

	if isValidLedgerEvent(stripeEvent) {
		getBalanceTxnsErr := svc.wallet.GetNewBalanceTransaction(ctx, extractedStripeEvent)
		if getBalanceTxnsErr != nil {
			return httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when updating ledger and balances",
				fmt.Sprintf("Something went wrong when updating ledger and balances. Err: %v", getBalanceTxnsErr),
			)
		}
	}

	return nil
}

func isValidLedgerEvent(event stripe.Event) bool {
	return event.Type == "charge.captured" || event.Type == "invoice.paid"
}

func isValidPendingChargeEvent(event stripe.Event) bool {
	return event.Type == "charge.expired"
}

func isValidSubscriptionEvent(event stripe.Event) bool {
	subscriptionID, isSubscriptionIncluded := event.Data.Object["subscription"]
	return (event.Type == "invoice.payment_failed" || event.Type == "invoice.paid") &&
		isSubscriptionIncluded &&
		subscriptionID != nil
}
