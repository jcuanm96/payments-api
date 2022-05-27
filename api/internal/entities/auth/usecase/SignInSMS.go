package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	twilio "github.com/kevinburke/twilio-go"
)

const signinErr = "Something went wrong when trying to sign in. Please try again."

func (svc *usecase) SignInSMS(ctx context.Context, req request.SignInSMS) (*response.AuthSuccess, error) {
	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signinErr,
			fmt.Sprintf("Could not begin transaction. Error: %v", txErr),
		)
	}

	formattedNumber, formatErr := req.FormattedPhonenumber()
	if formatErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signinErr,
			fmt.Sprintf("Could not format phone number %s. Err: %v", req.Number, formatErr),
		)
	}

	// Ignore Twilio code check if we're testing
	phonenumberCheckStatus := twilio.CheckPhoneNumber{
		Status: "approved",
	}

	// Ignore Twilio code check if we're testing
	if appconfig.Config.Gcloud.Project == "vama-prod" {
		status, twilioCheckErr := svc.twilio.Check(ctx, appconfig.Config.Twilio.Verify, url.Values{"To": []string{formattedNumber}, "Code": []string{req.Code}})
		if twilioCheckErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				401,
				http.StatusUnauthorized,
				"Invalid verification code.",
				fmt.Sprintf("Invalid verification code. Err: %v", twilioCheckErr),
			)
		}
		phonenumberCheckStatus = *status
	}

	if phonenumberCheckStatus.Status != "approved" {
		return nil, httperr.NewCtx(
			ctx,
			401,
			http.StatusUnauthorized,
			"Please enter a valid verification code.",
			"Please enter a valid verification code.",
		)
	}
	user, getUserByPhoneErr := svc.user.GetUserByPhone(ctx, tx, formattedNumber)
	if getUserByPhoneErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signinErr,
			fmt.Sprintf("Error retrieving user with phone number %s from db. Err: %v", formattedNumber, getUserByPhoneErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			fmt.Sprintf("No user associated with phone number %s", formattedNumber),
			fmt.Sprintf("No user associated with phone number %s", formattedNumber),
		)
	}

	// Defer a rollback unless we set commit to true
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	res, authUserErr := svc.AuthUser(ctx, tx, user.UUID, user.ID)

	if authUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signinErr,
			fmt.Sprintf("Something went wrong when authenticating phone %s. Err: %v", formattedNumber, authUserErr),
		)
	}

	sendbirdUser, getSendbirdUserErr := svc.sendbirdClient.GetUser(user.ID)
	if getSendbirdUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signinErr,
			fmt.Sprintf("Error getting sendbird user ID %d: %v", user.ID, getSendbirdUserErr),
		)
	}

	res.Credentials.SendBirdAccessToken = sendbirdUser.AccessToken
	res.User = *user
	commit = true
	return res, nil
}
