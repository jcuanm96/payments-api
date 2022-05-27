package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/jackc/pgx/v4"
	"github.com/stripe/stripe-go/v72"
)

const upsertGoatSubscriptionInfoErr = "Something went wrong when updating subscription price."

func (svc *usecase) UpsertGoatSubscriptionInfoTx(ctx context.Context, tx pgx.Tx, goatUser *response.User, priceInSmallestDenom int64, currency string) error {
	const tierName = constants.DEFAULT_TIER_NAME
	stripeProductName := subscription.FormatSubscriptionProduct(goatUser.Username, tierName, priceInSmallestDenom, currency)
	stripeProductParams := &stripe.ProductParams{
		Name: stripe.String(stripeProductName),
	}
	stripeProduct, createStripeProductErr := svc.stripeClient.Products.New(stripeProductParams)
	if createStripeProductErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			upsertGoatSubscriptionInfoErr,
			fmt.Sprintf("Could not create new Stripe product for creator %d: %v", goatUser.ID, createStripeProductErr),
		)
	}

	subscriptionErr := svc.repo.UpsertGoatSubscriptionInfo(ctx, tx, goatUser.ID, tierName, priceInSmallestDenom, currency, stripeProduct.ID)
	if subscriptionErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			upsertGoatSubscriptionInfoErr,
			fmt.Sprintf("Could not update tiers database: %v", subscriptionErr),
		)
	}
	return nil
}

func (svc *usecase) UpsertGoatSubscriptionInfo(ctx context.Context, priceInSmallestDenom int64, currency string) error {
	currUser, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			upsertGoatSubscriptionInfoErr,
			fmt.Sprintf("Could not find user in the current context: %v", getUserErr),
		)
	} else if currUser == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			upsertGoatSubscriptionInfoErr,
			"currUser was nil in UpsertGoatSubscriptionInfo",
		)
	} else if currUser.Type != "GOAT" {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can update subscription prices. Interested?  Sign up to be a creator now!",
			"Only creators can update subscription prices.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			upsertGoatSubscriptionInfoErr,
			fmt.Sprintf("Error starting txn: %v", txErr),
		)
	}
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	upsertSubscriptionHttpErr := svc.UpsertGoatSubscriptionInfoTx(ctx, tx, currUser, priceInSmallestDenom, currency)
	if upsertSubscriptionHttpErr != nil {
		return upsertSubscriptionHttpErr
	}

	commit = true
	return nil
}
