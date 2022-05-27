package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

const makePaymentIntentDefaultErr = "Something went wrong making transaction."

func (svc *usecase) MakePaymentIntent(ctx context.Context, req request.MakePaymentIntent) error {
	// Validating customer information
	currentUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			makePaymentIntentDefaultErr,
			fmt.Sprintf("Error occurred when retrieving user from db. Err: %v", getCurrUserErr),
		)
	} else if currentUser == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			makePaymentIntentDefaultErr,
			"Current user nil for MakePaymentIntent",
		)
	} else if currentUser.StripeID == nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Error making transaction",
			fmt.Sprintf("The Stripe ID for user %d does not exist", currentUser.ID),
		)
	} else if currentUser.Email == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"Please set your email in order to make a transaction",
			"Please set your email in order to make a transaction",
		)
	}

	currUserCustomerObj, getCustomerErr := svc.stripeClient.Customers.Get(*currentUser.StripeID, nil)
	if getCustomerErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			makePaymentIntentDefaultErr,
			fmt.Sprintf("Error getting Stripe customer info from user %d . Err: %v", currentUser.ID, getCustomerErr),
		)
	}

	if currUserCustomerObj.DefaultSource == nil || currUserCustomerObj.DefaultSource.ID == "" {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"No default payment method was found",
			"No default payment method was found",
		)
	}

	// Validating provider information
	providerUser, getUserByIDErr := svc.user.GetUserByID(ctx, req.ProviderUserID)

	if getUserByIDErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			makePaymentIntentDefaultErr,
			fmt.Sprintf("Error occurred when retrieving provider user %d from db. Err: %v", req.ProviderUserID, getUserByIDErr),
		)
	} else if providerUser == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"The provider you are trying to pay does not exist",
			"The provider you are trying to pay does not exist",
		)
	} else if providerUser.Type != "GOAT" {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can receive payments",
			"Only creators can receive payments",
		)
	} else if providerUser.ID == currentUser.ID {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"A user cannot pay themselves",
			"A user cannot pay themselves",
		)
	}

	captureMethod := stripe.PaymentIntentCaptureMethodManual // Will wait for our server to confirm the transfer of funds
	if req.AutoCapture {
		captureMethod = stripe.PaymentIntentCaptureMethodAutomatic
	}

	// Make the payment intent
	paymentIntentParams := &stripe.PaymentIntentParams{
		Confirm:       stripe.Bool(true),
		Customer:      currentUser.StripeID,
		CaptureMethod: stripe.String(string(captureMethod)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		PaymentMethod: stripe.String(currUserCustomerObj.DefaultSource.ID),
		Amount:        stripe.Int64(req.AmountInSmallestDenom),
		Currency:      stripe.String(req.Currency),
	}
	paymentIntentParams.AddMetadata("providerUserID", strconv.Itoa(providerUser.ID))
	paymentIntentParams.AddMetadata("customerUserID", strconv.Itoa(currentUser.ID))

	paymentIntent, createPaymentIntentErr := svc.stripeClient.PaymentIntents.New(paymentIntentParams)

	if createPaymentIntentErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			makePaymentIntentDefaultErr,
			fmt.Sprintf("Error making payment intent for customer %d and provider %d: %v", currentUser.ID, req.ProviderUserID, createPaymentIntentErr),
		)
	}

	if paymentIntent == nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			makePaymentIntentDefaultErr,
			fmt.Sprintf("Nil payment intent for customer %d and provider %d.", currentUser.ID, req.ProviderUserID),
		)
	}

	if !req.AutoCapture {
		insertChargeErr := svc.repo.InsertPendingPayment(ctx, svc.repo.MasterNode(), *paymentIntent, currentUser.ID, providerUser.ID)
		if insertChargeErr != nil {
			return httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				makePaymentIntentDefaultErr,
				fmt.Sprintf("Error inserting pending payment %s: %v", paymentIntent.ID, insertChargeErr),
			)
		}
	}
	return nil

}
