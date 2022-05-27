package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/kevinburke/twilio-go"
	"github.com/stripe/stripe-go/v72"
)

const errChangingPhoneNumber = "Something went wrong when trying to change your phone number. Please try again."

func (svc *usecase) UpdatePhone(ctx context.Context, req request.UpdatePhone) (*response.UpdatePhone, error) {
	user, userErr := svc.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errChangingPhoneNumber,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	} else if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errChangingPhoneNumber,
			"user was nil in UpdatePhone",
		)
	}
	phoneNumberCheckStatus := twilio.CheckPhoneNumber{
		Status: "approved",
	}
	formattedNumber, formatErr := req.FormattedPhonenumber()
	if formatErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errChangingPhoneNumber,
			fmt.Sprintf("Could not format phone number %s: %v", req.Number, formatErr),
		)
	}

	taken, checkPhoneErr := svc.CheckPhoneAlreadyExists(ctx, formattedNumber)
	if checkPhoneErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errChangingPhoneNumber,
			fmt.Sprintf("Couldn't verify if phone was already taken: %v", checkPhoneErr),
		)
	} else if taken {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"The phone number you entered has already been taken.  Try logging in instead?",
			"The phone number you entered has already been taken.",
		)
	}

	status, twilioErr := svc.twilio.Check(ctx, svc.twilioVerifyID, url.Values{"To": []string{formattedNumber}, "Code": []string{req.Code}})
	if twilioErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Invalid verification code.",
			fmt.Sprintf("Invalid verification code: %v", twilioErr),
		)
	}
	phoneNumberCheckStatus = *status

	if phoneNumberCheckStatus.Status != "approved" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Invalid verification code.",
			"Invalid verification code.",
		)
	}

	res := response.UpdatePhone{}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errChangingPhoneNumber,
			fmt.Sprintf("Could not begin transaction: %v", txErr),
		)
	}
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	updateErr := svc.repo.UpdatePhone(ctx, tx, user.ID, formattedNumber)
	if updateErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errChangingPhoneNumber,
			fmt.Sprintf("Error updating phone in DB: %v", updateErr),
		)
	}

	stripeParams := &stripe.CustomerParams{
		Phone: &formattedNumber,
	}

	// Update Stripe params
	if user.StripeID != nil {
		_, stripeErr := svc.stripeClient.Customers.Update(
			*user.StripeID,
			stripeParams,
		)

		if stripeErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				errChangingPhoneNumber,
				fmt.Sprintf("Error updating Stripe phone number for user: %d: %v", user.ID, stripeErr),
			)
		}
	}

	commit = true
	return &res, nil
}
