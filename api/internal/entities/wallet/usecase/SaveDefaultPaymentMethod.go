package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

type StripeTokenError struct {
	Message string `json:"message"`
}

const defaultPaymentMethodErr = "There was a problem saving your payment information. Please try again."

func (svc *usecase) SaveDefaultPaymentMethod(ctx context.Context, req request.DefaultPaymentMethod) (*response.DefaultPaymentMethod, error) {
	currentUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			defaultPaymentMethodErr,
			fmt.Sprintf("Error occurred when retrieving user %d from db. Err: %v", currentUser.ID, getCurrUserErr),
		)
	} else if currentUser.StripeID == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			defaultPaymentMethodErr,
			fmt.Sprintf("The Stripe ID for user %d does not exist", currentUser.ID),
		)
	} else if currentUser.Email == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"Please set your email in order to save card information.",
			"Please set your email in order to save card information.",
		)
	}

	expYear, AtoiErr := strconv.Atoi(req.ExpYear)
	if AtoiErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Please enter a valid expiration year.",
			fmt.Sprintf("Unable to convert cc expiration year to int. Err: %v", AtoiErr),
		)
	}

	expMonth, expMonthAtoiErr := strconv.Atoi(req.ExpMonth)
	if expMonthAtoiErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Please enter a valid expiration month.",
			fmt.Sprintf("Unable to convert credit card expiration month to valid int. Err: %v", expMonthAtoiErr),
		)
	} else if expMonth < 1 || expMonth > 12 {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Please enter a valid expiration month.",
			"Expiration date must be between 1 and 12 inclusive",
		)
	}

	currentYear, currentMonth, _ := time.Now().Date()
	if expYear < currentYear || (expYear == currentYear && expMonth < int(currentMonth)) {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Credit card is expired.",
			"Credit card is expired.",
		)
	}

	// Get credit card token
	tokenParams := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   stripe.String(req.Number),
			ExpMonth: stripe.String(req.ExpMonth),
			ExpYear:  stripe.String(req.ExpYear),
			CVC:      stripe.String(req.CVC),
		},
	}

	creditCardToken, newTokenErr := svc.stripeClient.Tokens.New(tokenParams)
	if newTokenErr != nil {
		var stripeErr StripeTokenError
		json.Unmarshal([]byte(newTokenErr.Error()), &stripeErr)
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Something went wrong saving your payment information: %s", stripeErr.Message),
			fmt.Sprintf("Something went wrong when saving user %d card info. TokenErr: %v. StripeErr: %v", currentUser.ID, newTokenErr, stripeErr.Message),
		)
	}

	params := &stripe.CustomerParams{}
	params.SetSource(creditCardToken.ID)
	cus, updateCustomerErr := svc.stripeClient.Customers.Update(*currentUser.StripeID, params)
	if updateCustomerErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			defaultPaymentMethodErr,
			fmt.Sprintf("Something went wrong updating Stripe customer for user %d. Err: %v", currentUser.ID, updateCustomerErr),
		)
	} else if cus.DefaultSource == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			defaultPaymentMethodErr,
			fmt.Sprintf("Something went wrong the default source returned by stripe is nil when updating Stripe customer for user %d.", currentUser.ID),
		)
	}

	res := response.DefaultPaymentMethod{}
	return &res, nil
}
