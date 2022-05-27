package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/nyaruka/phonenumbers"
	"github.com/stripe/stripe-go/v72"
)

const verifyProviderAccountErr = "Something went wrong when trying to create a payment account."

func (svc *usecase) VerifyProviderAccount(ctx context.Context, req request.VerifyProviderAccount) (*response.VerifyProviderAccount, error) {
	stripeAccountID := req.StripeAccountID
	phoneNumber := getFormattedPhonenumber(req.Number, req.CountryCode)
	if stripeAccountID == nil {
		user, getUserErr := svc.user.GetUserByPhone(ctx, svc.repo.MasterNode(), phoneNumber)
		if getUserErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				verifyProviderAccountErr,
				fmt.Sprintf("Error could not get user for phone number %s. Error: %v", phoneNumber, getUserErr),
			)
		} else if user == nil {
			return nil, httperr.NewCtx(
				ctx,
				404,
				http.StatusNotFound,
				verifyProviderAccountErr,
				"user returned nil for VerifyProviderAccount",
			)
		}

		if user.StripeAccountID == nil {
			res := response.VerifyProviderAccount{
				StripeAccountID:     "",
				HasDetailsSubmitted: false,
			}
			return &res, nil
		}

		stripeAccountID = user.StripeAccountID
	}

	stripeAccount, getStripeAccountErr := svc.stripeClient.Account.GetByID(
		*stripeAccountID,
		nil,
	)

	if getStripeAccountErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			verifyProviderAccountErr,
			fmt.Sprintf("Error could not get Stripe account information for account_id %s. Error: %v", *stripeAccountID, getStripeAccountErr),
		)
	} else if stripeAccount == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			verifyProviderAccountErr,
			"Error the Stripe account object was returned as nil.",
		)
	}

	res := response.VerifyProviderAccount{
		StripeAccountID:     *stripeAccountID,
		HasDetailsSubmitted: isStripeBankVerificationComplete(stripeAccount),
	}

	return &res, nil
}

// The user is verified if they have done the following:
// 1. Clicked submit on the form
// 2. Provided the necessary documents (eg. id, bank account information)
// 3. Nothing is pending
func isStripeBankVerificationComplete(stripeAccount *stripe.Account) bool {
	return stripeAccount.DetailsSubmitted &&
		stripeAccount.Requirements.CurrentDeadline == 0 &&
		stripeAccount.Requirements.DisabledReason == "" &&
		len(stripeAccount.Requirements.CurrentlyDue) == 0 &&
		len(stripeAccount.Requirements.Errors) == 0 &&
		len(stripeAccount.Requirements.EventuallyDue) == 0 &&
		len(stripeAccount.Requirements.PastDue) == 0 &&
		len(stripeAccount.Requirements.PendingVerification) == 0
}

func getFormattedPhonenumber(number string, countryCode string) string {
	num, err := phonenumbers.Parse(number, countryCode)
	if err != nil {
		return ""
	}
	return phonenumbers.Format(num, phonenumbers.E164)
}
