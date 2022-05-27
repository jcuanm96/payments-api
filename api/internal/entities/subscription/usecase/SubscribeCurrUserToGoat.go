package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

const subscribeCurrUserToGoatErr = "Something went wrong when trying to subscribe to creator."

func (svc *usecase) SubscribeCurrUserToGoat(ctx context.Context, req request.SubscribeCurrUserToGoat) (*response.UserSubscription, error) {
	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("Something went wrong when creating transaction: %v", txErr),
		)
	}
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	var usersWg sync.WaitGroup

	var currUser *response.User
	var getCurrentUserErr error
	usersWg.Add(1)
	go func() {
		defer usersWg.Done()
		currUser, getCurrentUserErr = svc.user.GetCurrentUser(ctx)
	}()

	var goatUser *response.User
	var getGoatUserErr error
	usersWg.Add(1)
	go func() {
		defer usersWg.Done()
		goatUser, getGoatUserErr = svc.user.GetUserByID(ctx, req.GoatUserID)
	}()

	usersWg.Wait()

	if getCurrentUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("Could not find user in the current context: %v", getCurrentUserErr),
		)
	} else if currUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			subscribeCurrUserToGoatErr,
			"currUser was nil in SubscribeCurrUserToGoat.",
		)
	} else if currUser.StripeID == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("User %d does not have a Stripe customer ID", currUser.ID),
		)
	} else if currUser.Email == nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You must set your email before being able to subscribe.",
			"You must set your email before being able to subscribe.",
		)
	}

	if getGoatUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("Could not get creator user %d: %v", req.GoatUserID, getGoatUserErr),
		)
	} else if goatUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			subscribeCurrUserToGoatErr,
			"goatUser was nil in SubscribeCurrUserToGoat.",
		)
	} else if goatUser.Type != "GOAT" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"The person you're trying to subscribe to is not a creator.",
			"The person you're trying to subscribe to is not a creator.",
		)
	} else if goatUser.ID == currUser.ID {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You cannot subscribe to yourself.",
			"You cannot subscribe to yourself.",
		)
	}

	var subscriptionInfoWg sync.WaitGroup

	var currUserStripeCustomer *stripe.Customer
	var getStripeCustomerErr error
	subscriptionInfoWg.Add(1)
	go func() {
		defer subscriptionInfoWg.Done()
		currUserStripeCustomer, getStripeCustomerErr = svc.stripeClient.Customers.Get(*currUser.StripeID, nil)
	}()

	var oldSubscription *response.UserSubscription
	var getSubscriptionByGoatIDErr error
	subscriptionInfoWg.Add(1)
	go func() {
		defer subscriptionInfoWg.Done()
		// Check if the user is already subscribed to the creator
		oldSubscription, getSubscriptionByGoatIDErr = svc.repo.GetUserSubscriptionByGoatID(ctx, currUser.ID, goatUser.ID)
	}()

	subscriptionTier, getGoatSubTierErr := svc.repo.GetGoatSubscriptionInfo(ctx, goatUser.ID)
	if getGoatSubTierErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("Could not get subscription tier from DB for user %d: %v", goatUser.ID, getGoatSubTierErr),
		)
	} else if subscriptionTier == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("Subscription tier for user %d does not exist.", goatUser.ID),
		)
	}

	subscriptionInfoWg.Wait()

	// Check if user has a default payment method set
	if getStripeCustomerErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("Error getting Stripe customer for user: %d: %v", currUser.ID, getStripeCustomerErr),
		)
	} else if currUserStripeCustomer.DefaultSource == nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You must set a default payment method in order to subscribe.",
			"You must set a default payment method in order to subscribe.",
		)
	}

	if getSubscriptionByGoatIDErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("Failed when checking if user %d is already subscribed to creator %d: %v", currUser.ID, goatUser.ID, getSubscriptionByGoatIDErr),
		)
	} else if isSubscribed(oldSubscription) && oldSubscription.IsRenewing {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You are already subscribed to this creator.",
			"You are already subscribed to this creator.",
		)
	}

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(*currUser.StripeID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				PriceData: &stripe.SubscriptionItemPriceDataParams{
					Currency: stripe.String(subscriptionTier.Currency),
					Product:  stripe.String(subscriptionTier.StripeProductID),
					Recurring: &stripe.SubscriptionItemPriceDataRecurringParams{
						Interval: stripe.String("month"),
					},
					UnitAmount: stripe.Int64(subscriptionTier.PriceInSmallestDenom),
				},
			},
		},
		PaymentBehavior: stripe.String("error_if_incomplete"),
	}
	params.AddExpand("latest_invoice.payment_intent")
	params.AddMetadata("customerUserID", strconv.Itoa(currUser.ID))
	params.AddMetadata("providerUserID", strconv.Itoa(req.GoatUserID))
	params.AddMetadata("tierID", strconv.Itoa(subscriptionTier.ID))

	// User unsubscribed and resubscribed before the current period ended,
	// so we don't charge them for the new subscription until the current period ends.
	// We give them a temporary 'free trial'. An `invoice.paid` event is triggered
	// so the status will be updated from the Stripe subscription webhook.
	if isSubscribed(oldSubscription) && !oldSubscription.IsRenewing {
		params.TrialEnd = stripe.Int64(oldSubscription.CurrentPeriodEnd.Unix())
	}

	var newCurrentPeriodEnd time.Time
	if params.TrialEnd != nil {
		newCurrentPeriodEnd = oldSubscription.CurrentPeriodEnd
	} else {
		years := 0
		months := 1
		days := 0
		newCurrentPeriodEnd = time.Now().AddDate(years, months, days)
	}

	newSubscription := response.UserSubscription{
		UserID:               currUser.ID,
		StripeSubscriptionID: "",
		GoatUser:             *goatUser,
		TierID:               subscriptionTier.ID,
		IsRenewing:           true,
		CurrentPeriodEnd:     newCurrentPeriodEnd,
	}

	upsertSubscriptionErr := svc.UpsertUserSubscription(ctx, tx, newSubscription)
	if upsertSubscriptionErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("Error upserting subscription: %v", upsertSubscriptionErr),
		)
	}

	stripeSubscription, createSubscriptionErr := svc.stripeClient.Subscriptions.New(params)
	if createSubscriptionErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("Error creating a Stripe subscription object for user %d, creator %d. Err: %v", currUser.ID, goatUser.ID, createSubscriptionErr),
		)
	} else if stripeSubscription.LatestInvoice == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribeCurrUserToGoatErr,
			fmt.Sprintf("Error the Stripe subscription %s has a nil invoice for provider %d and tier id %d.", stripeSubscription.ID, goatUser.ID, subscriptionTier.ID),
		)
	}

	commit = true
	return &newSubscription, nil
}

func isSubscribed(userSubscription *response.UserSubscription) bool {
	return userSubscription != nil && time.Now().Before(userSubscription.CurrentPeriodEnd)
}
