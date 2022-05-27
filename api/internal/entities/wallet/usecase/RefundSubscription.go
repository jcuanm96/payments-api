package service

import (
	"context"

	"github.com/stripe/stripe-go/v72"
)

func (svc *usecase) RefundSubscription(ctx context.Context, paymentIntentID string) (*stripe.Refund, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
	}

	refund, refundErr := svc.stripeClient.Refunds.New(params)
	return refund, refundErr
}
