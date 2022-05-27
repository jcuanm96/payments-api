package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

const getMyPaymentMethodsErr = "Something went wrong when getting payment methods."

// GetMyPaymentMethods retrieves an array of payment methods.
// All users have a default payment method (card).
// Only users with a stripe account (creators) have an external payment method (such as a bank account).
// This function makes the card request and the external method request (if applicable) in parallel.
// Once both are done, we process the errors/results and aggregate into an array.
func (svc *usecase) GetMyPaymentMethods(ctx context.Context) (*response.GetMyPaymentMethods, error) {
	currentUser, currentUserErr := svc.user.GetCurrentUser(ctx)
	if currentUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getMyPaymentMethodsErr,
			fmt.Sprintf("Something went wrong getting current user: %v", currentUserErr),
		)
	}
	if currentUser == nil {
		vlog.Errorf(ctx, "Current user nil for GetMyPaymentMethods")
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			getMyPaymentMethodsErr,
			"Current user nil for GetMyPaymentMethods",
		)
	}

	var wg sync.WaitGroup

	// Get the default payment method (card) for all users.
	var cus *stripe.Customer
	var getCustomerErr error
	var card *stripe.Card
	var getCardErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		cus, getCustomerErr = svc.stripeClient.Customers.Get(*currentUser.StripeID, nil)
		// The error is non-nil, so return this goroutine and the following handling
		// code will log and return the appropriate errors.
		if getCustomerErr != nil {
			return
		} else if cus.DefaultSource == nil {
			return
		}
		cardParams := &stripe.CardParams{
			Customer: stripe.String(cus.ID),
		}

		card, getCardErr = svc.stripeClient.Cards.Get(
			cus.DefaultSource.ID,
			cardParams,
		)
	}()

	var banks []response.PaymentMethod
	var getBanksErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		banks, getBanksErr = svc.repo.GetUserBankInfo(ctx, currentUser.ID)
	}()

	wg.Wait()

	// Process card result/errors for all users
	if getCustomerErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getMyPaymentMethodsErr,
			fmt.Sprintf("Something went wrong getting Stripe customer for user: %d, err: %v", currentUser.ID, getCustomerErr),
		)
	}

	if getCardErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getMyPaymentMethodsErr,
			fmt.Sprintf("Something went wrong getting Stripe card for user: %d, err: %v", currentUser.ID, getCardErr),
		)
	}

	if getBanksErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getMyPaymentMethodsErr,
			fmt.Sprintf("Something went wrong getting banks for user %d. Err: %v", currentUser.ID, getBanksErr),
		)
	} else if banks == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getMyPaymentMethodsErr,
			"Error banks slice came back as nil",
		)
	}

	paymentMethods := []response.PaymentMethod{}
	if card != nil {
		paymentMethod := response.PaymentMethod{
			Card: &response.Card{
				Brand:    fmt.Sprintf("%v", card.Brand),
				ExpMonth: card.ExpMonth,
				ExpYear:  card.ExpYear,
				Funding:  fmt.Sprintf("%v", card.Funding),
				Last4:    card.Last4,
				Country:  card.Country,
				CVCCheck: fmt.Sprintf("%v", card.CVCCheck),
				Name:     card.Name,
			},
		}
		paymentMethods = append(paymentMethods, paymentMethod)
	}

	paymentMethods = append(paymentMethods, banks...)

	res := response.GetMyPaymentMethods{
		PaymentMethods: paymentMethods,
	}

	return &res, nil
}
