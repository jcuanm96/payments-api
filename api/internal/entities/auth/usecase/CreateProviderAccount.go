package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

const createProviderAccountErr = "Something went wrong when trying to create a payment account."

func (svc *usecase) CreateProviderAccount(ctx context.Context, req request.CreateProviderAccount) (*response.CreateProviderAccount, error) {
	user, getUserErr := svc.user.GetUserByEmail(ctx, svc.repo.MasterNode(), req.Email)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			createProviderAccountErr,
			fmt.Sprintf("Error occurred when getting the user by their email %s. Err: %v", req.Email, getUserErr),
		)
	}
	currStripeAccountID := req.StripeAccountID

	if user != nil && user.StripeAccountID != nil {
		currStripeAccountID = user.StripeAccountID
	}

	// Generate a new stripeAccountID if none is passed into the request
	// This ensures we don't make multiple accounts for the same user
	if req.StripeAccountID == nil {
		accountParams := &stripe.AccountParams{
			Settings: &stripe.AccountSettingsParams{
				Payouts: &stripe.AccountSettingsPayoutsParams{
					Schedule: &stripe.PayoutScheduleParams{
						Interval: stripe.String("manual"),
					},
				},
			},
			Capabilities: &stripe.AccountCapabilitiesParams{
				CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
					Requested: stripe.Bool(true),
				},
				Transfers: &stripe.AccountCapabilitiesTransfersParams{
					Requested: stripe.Bool(true),
				},
			},
			Country:      stripe.String("US"),
			BusinessType: stripe.String(fmt.Sprintf("%v", stripe.AccountBusinessTypeIndividual)),
			Email:        stripe.String(req.Email),
			Type:         stripe.String("express"),
		}
		stripeAccount, createAccountErr := svc.stripeClient.Account.New(accountParams)
		currStripeAccountID = &stripeAccount.ID

		if createAccountErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				createProviderAccountErr,
				fmt.Sprintf("Error occurred when creating Stripe Connect account for email %s. Err: %v", req.Email, createAccountErr),
			)
		}

		upsertErr := svc.user.UpsertStripeAccountIDByEmail(ctx, req.Email, *currStripeAccountID)
		if upsertErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				createProviderAccountErr,
				fmt.Sprintf("Error upserting the user's stripe account id. Err: %v", upsertErr),
			)
		}
	}

	// Create Stripe Account link to verify creator's payment background information
	refreshLink := fmt.Sprintf("%s/goat-onboarding/refresh?account_id=%s&env=%s&email=%s", constants.VAMA_WEB_BASE_URL, *currStripeAccountID, appconfig.Config.Gcloud.Project, req.Email)
	returnLink := fmt.Sprintf("%s/goat-onboarding/return?account_id=%s&env=%s", constants.VAMA_WEB_BASE_URL, *currStripeAccountID, appconfig.Config.Gcloud.Project)
	accountLinkParams := &stripe.AccountLinkParams{
		Account:    stripe.String(*currStripeAccountID),
		RefreshURL: stripe.String(refreshLink),
		ReturnURL:  stripe.String(returnLink),
		Type:       stripe.String("account_onboarding"),
	}

	accountLink, createAccountErr := svc.stripeClient.AccountLinks.New(accountLinkParams)

	if createAccountErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			createProviderAccountErr,
			fmt.Sprintf("Error occurred when creating Stripe Connect account link for email %s. Err: %v", req.Email, createAccountErr),
		)
	}

	res := response.CreateProviderAccount{
		URL:       accountLink.URL,
		CreatedAt: accountLink.Created,
		ExpiresAt: accountLink.ExpiresAt,
		Object:    accountLink.Object,
	}

	return &res, nil
}
