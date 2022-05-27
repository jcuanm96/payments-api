package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

const confirmPaymentErr = "There was a problem making the transaction."

func (svc *usecase) ConfirmPaymentIntent(ctx context.Context, req request.ConfirmPaymentIntent) (*response.ConfirmPaymentIntent, error) {
	// Validating customer information
	currentUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			confirmPaymentErr,
			fmt.Sprintf("Error occurred when retrieving user from db. Err: %s", getCurrUserErr),
		)
	} else if currentUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusInternalServerError,
			confirmPaymentErr,
			"Current user is nil for ConfirmPaymentIntent.",
		)
	} else if currentUser.Type != "GOAT" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can confirm a payment",
			fmt.Sprintf("User %d not of type GOAT, was type %s", currentUser.ID, currentUser.Type),
		)
	}

	paymentIntentID, getPaymentIntentErr := svc.repo.GetPaymentIntent(ctx, req.CustomerUserID, currentUser.ID)
	if getPaymentIntentErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"There was a problem completing the transaction",
			fmt.Sprintf("Something went wrong retrieving payment intent for db for users %d and %d. Err: %v", req.CustomerUserID, currentUser.ID, getPaymentIntentErr),
		)
	}

	if paymentIntentID == "" {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			confirmPaymentErr,
			"No payment intent to confirm was found",
		)
	}

	paymentIntentCaptureParams := &stripe.PaymentIntentCaptureParams{}
	paymentIntent, captureFundsErr := svc.stripeClient.PaymentIntents.Capture(paymentIntentID, paymentIntentCaptureParams)

	if captureFundsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			confirmPaymentErr,
			fmt.Sprintf("Something went wrong confirming payment intent for users %d and %d. Err: %v", req.CustomerUserID, currentUser.ID, captureFundsErr),
		)
	} else if paymentIntent.Charges == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			confirmPaymentErr,
			fmt.Sprintf("Something went wrong the confirmed payment intent came back nil for user %d.", currentUser.ID),
		)
	} else if len(paymentIntent.Charges.Data) == 0 {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			confirmPaymentErr,
			fmt.Sprintf("Something went wrong the confirmed payment intent has no charges for user %d.", currentUser.ID),
		)
	}

	res := response.ConfirmPaymentIntent{
		ChargeID: paymentIntent.Charges.Data[0].ID,
	}

	return &res, nil
}
