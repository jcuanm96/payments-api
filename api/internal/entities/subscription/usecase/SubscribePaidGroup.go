package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/jackc/pgx/v4"
	"github.com/stripe/stripe-go/v72"
)

const subscribePaidGroupErr = "Something went wrong when trying to subscribe to paid group chat."

func (svc *usecase) SubscribePaidGroup(ctx context.Context, tx pgx.Tx, channelID string, user *response.User) (*response.PaidGroupChatSubscription, error) {
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			subscribePaidGroupErr,
			"user was nil in SubscribePaidGroup.",
		)
	} else if user.StripeID == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribePaidGroupErr,
			fmt.Sprintf("User %d does not have a Stripe customer ID", user.ID),
		)
	} else if user.Email == nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You must set your email before being able to join paid group chat.",
		)
	}

	var wg sync.WaitGroup

	var userStripeCustomer *stripe.Customer
	var getStripeCustomerErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		userStripeCustomer, getStripeCustomerErr = svc.stripeClient.Customers.Get(*user.StripeID, nil)
	}()

	var oldSubscription *response.PaidGroupChatSubscription
	var getSubscriptionErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		oldSubscription, getSubscriptionErr = svc.repo.GetPaidGroupSubscription(ctx, tx, user.ID, channelID)
	}()

	var paidGroupInfo *subscription.PaidGroupChatInfo
	var getPaidGroupInfoErr error
	var goatUser *response.User
	var getGoatUserErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		paidGroupInfo, getPaidGroupInfoErr = svc.repo.GetPaidGroupProductInfo(ctx, svc.repo.MasterNode(), channelID)
		if getPaidGroupInfoErr != nil {
			return
		}
		goatUser, getGoatUserErr = svc.user.GetUserByID(ctx, paidGroupInfo.GoatID)
	}()

	wg.Wait()

	if getStripeCustomerErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribePaidGroupErr,
			fmt.Sprintf("Error getting stripe customer ID %s: %v", *user.StripeID, getStripeCustomerErr),
		)
	} else if userStripeCustomer.DefaultSource == nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You must set a default payment method in order to join paid group chat.",
		)
	}

	if getSubscriptionErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribePaidGroupErr,
			fmt.Sprintf("Error getting existing paid group subscription: %v", getSubscriptionErr),
		)
	} else if isSubscribedToPaidGroup(oldSubscription) && oldSubscription.IsRenewing {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You are already subscribed to this paid group chat.",
		)
	}

	if getPaidGroupInfoErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribePaidGroupErr,
			fmt.Sprintf("Error getting paid group subscription info for channel %s: %v", channelID, getPaidGroupInfoErr),
		)
	}

	if getGoatUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribePaidGroupErr,
			fmt.Sprintf("Error getting creator user %d: %v", paidGroupInfo.GoatID, getGoatUserErr),
		)
	} else if goatUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribePaidGroupErr,
			fmt.Sprintf("Goat user %d was nil", paidGroupInfo.GoatID),
		)
	}

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(*user.StripeID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				PriceData: &stripe.SubscriptionItemPriceDataParams{
					Currency: stripe.String(paidGroupInfo.Currency),
					Product:  stripe.String(paidGroupInfo.StripeProductID),
					Recurring: &stripe.SubscriptionItemPriceDataRecurringParams{
						Interval: stripe.String("month"),
					},
					UnitAmount: stripe.Int64(paidGroupInfo.PriceInSmallestDenom),
				},
			},
		},
		PaymentBehavior: stripe.String("error_if_incomplete"),
	}
	params.AddExpand("latest_invoice.payment_intent")
	params.AddMetadata("customerUserID", strconv.Itoa(user.ID))
	params.AddMetadata("providerUserID", strconv.Itoa(paidGroupInfo.GoatID))
	params.AddMetadata("sendbirdChannelID", channelID)

	// User unsubscribed and resubscribed before the current period ended,
	// so we don't charge them for the new subscription until the current period ends.
	// We give them a temporary 'free trial'. An `invoice.paid` event is triggered
	// so the status will be updated from the Stripe subscription webhook.
	if isSubscribedToPaidGroup(oldSubscription) && !oldSubscription.IsRenewing {
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

	newSubscription := &response.PaidGroupChatSubscription{
		UserID:               user.ID,
		StripeSubscriptionID: "",
		GoatUser:             *goatUser,
		ChannelID:            channelID,
		IsRenewing:           true,
		CurrentPeriodEnd:     newCurrentPeriodEnd,
	}

	upsertPaidGroupSubscriptionErr := svc.UpsertPaidGroupSubscription(ctx, tx, newSubscription)
	if upsertPaidGroupSubscriptionErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribePaidGroupErr,
			fmt.Sprintf("Error upserting paid group subscription: %v", upsertPaidGroupSubscriptionErr),
		)
	}

	stripeSubscription, createSubscriptionErr := svc.stripeClient.Subscriptions.New(params)
	if createSubscriptionErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribePaidGroupErr,
			fmt.Sprintf("Error creating a Stripe subscription object for user %d: %v", user.ID, createSubscriptionErr),
		)
	} else if stripeSubscription.LatestInvoice == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			subscribePaidGroupErr,
			fmt.Sprintf("Error the Stripe subscription %s has a nil invoice.", stripeSubscription.ID),
		)
	}

	newSubscription.CurrentPeriodEnd = time.Unix(stripeSubscription.CurrentPeriodEnd, 0)
	return newSubscription, nil
}

func isSubscribedToPaidGroup(userSubscription *response.PaidGroupChatSubscription) bool {
	return userSubscription != nil && time.Now().Before(userSubscription.CurrentPeriodEnd)
}
