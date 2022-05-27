package test_util

import (
	"context"
	"testing"
	"time"

	test_util "github.com/VamaSingapore/vama-api/cmd/test/util"
	"github.com/VamaSingapore/vama-api/internal/controller"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	chatrepo "github.com/VamaSingapore/vama-api/internal/entities/chat/repositories"
	subscriptionrepo "github.com/VamaSingapore/vama-api/internal/entities/subscription/repositories"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stripe/stripe-go/v72"
)

// Create a new paid group in the db that contains valid Stripe product
func FillBootstrapPaidGroupChatData(db *pgxpool.Pool, ctr *controller.Ctr, userID int, sendbirdChannelID string, isMemberLimitEnabled bool) {
	ctx := context.Background()

	stripeProduct, newProductErr := ctr.Stripe.Products.New(&stripe.ProductParams{
		Name: stripe.String("Test Product Name"),
	})

	if newProductErr != nil {
		vlog.Fatalf(ctx, "Error creating new test product: %v", newProductErr)
	}

	priceInSmallestDenom := 300
	currency := "usd"
	link := "fakeLink"
	memberLimit := 10
	metadata := []byte{}
	_, upsertErr := chatrepo.UpsertPaidGroupChat(
		ctx,
		userID,
		sendbirdChannelID,
		stripeProduct.ID,
		priceInSmallestDenom,
		currency,
		link,
		memberLimit,
		isMemberLimitEnabled,
		metadata,
		db)
	if upsertErr != nil {
		vlog.Fatalf(ctx, "Error upserting paid group chat: %v", upsertErr)
	}
}

func FillBootstrapPaidGroupChatWithMembers(t *testing.T, app *fiber.App, db *pgxpool.Pool, ctr *controller.Ctr) {
	// Create paid group
	priceInSmallestDenom := 300
	testGroupName := "Test Paid Group Chat 1"
	createPaidGroupParams := map[string]interface{}{
		"priceInSmallestDenom": priceInSmallestDenom,
		"currency":             "usd",
		"name":                 testGroupName,
		"linkSuffix":           "testpaidgroup1",
	}

	asUserCreatorID := 1
	createEndpoint := "/api/v1/chat/paid/group/create"
	createPaidGroupResp := &response.CreatePaidGroupChannel{}
	test_util.MakePostRequestAssert200(t, app, createEndpoint, createPaidGroupParams, &asUserCreatorID, createPaidGroupResp)

	// Have user join
	FillBootstrapCustomerUserSignupData(t, app)
	sendbirdChannelID := "sendbird_group_channel_123456789"
	joinPaidGroupParams := map[string]interface{}{
		"channelID": sendbirdChannelID,
	}

	asCustomerID := 2
	joinEndpoint := "/api/v1/chat/paid/group/join"
	joinPaidGroupResp := &response.PaidGroupChatSubscription{}
	test_util.MakePostRequestAssert200(t, app, joinEndpoint, joinPaidGroupParams, &asCustomerID, joinPaidGroupResp)

	// Emulate the webhook by manually updating the subscription
	ctx := context.Background()

	stripeProduct, newProductErr := ctr.Stripe.Products.New(&stripe.ProductParams{
		Name: stripe.String("Test Product Name"),
	})

	if newProductErr != nil {
		vlog.Fatalf(ctx, "Error creating new test product: %v", newProductErr)
	}

	tokenParams := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   stripe.String("4242424242424242"),
			ExpMonth: stripe.String("12"),
			ExpYear:  stripe.String("2030"),
			CVC:      stripe.String("123"),
		},
	}

	creditCardToken, newTokenErr := ctr.Stripe.Tokens.New(tokenParams)
	if newTokenErr != nil {
		vlog.Fatalf(ctx, "Error creating test Stripe credit card for paid group chat: %v", newTokenErr)
	}

	customer, newCustomerErr := ctr.Stripe.Customers.New(&stripe.CustomerParams{
		Source: &stripe.SourceParams{
			Token: stripe.String(creditCardToken.ID),
		},
	})

	if newCustomerErr != nil {
		vlog.Fatalf(ctx, "Error creating test Stripe customer for paid group chat: %v", newCustomerErr)
	}
	stripeSubscription, newSubscriptionErr := ctr.Stripe.Subscriptions.New(&stripe.SubscriptionParams{
		Customer: stripe.String(customer.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				PriceData: &stripe.SubscriptionItemPriceDataParams{
					Currency: stripe.String("usd"),
					Product:  stripe.String(stripeProduct.ID),
					Recurring: &stripe.SubscriptionItemPriceDataRecurringParams{
						Interval: stripe.String("month"),
					},
					UnitAmount: stripe.Int64(int64(priceInSmallestDenom)),
				},
			},
		},
		PaymentBehavior: stripe.String("error_if_incomplete"),
	})

	if newSubscriptionErr != nil {
		vlog.Fatalf(ctx, "Error creating new test subscription: %v", newSubscriptionErr)
	}
	newSubscription := &response.PaidGroupChatSubscription{
		StripeSubscriptionID: stripeSubscription.ID,
		CurrentPeriodEnd:     time.Now().AddDate(0, 1, 0), // Add one month
		UserID:               asCustomerID,
		GoatUser:             response.User{ID: asUserCreatorID},
		ChannelID:            sendbirdChannelID,
		IsRenewing:           true,
	}

	subscriptionrepo.UpsertPaidGroupSubscription(ctx, db, newSubscription)
}
