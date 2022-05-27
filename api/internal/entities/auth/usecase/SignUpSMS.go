package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	twilio "github.com/kevinburke/twilio-go"
	"github.com/stripe/stripe-go/v72"
)

const signupErr = "Something went wrong when trying to sign up. Please try again."

func (svc *usecase) SignupSMS(ctx context.Context, req request.SignupSMS) (*response.AuthSuccess, error) {
	formattedNumber, formatErr := req.FormattedPhonenumber()
	if formatErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Could not format phone number %s. Err: %v", req.Number, formatErr),
		)
	}
	// Ignore Twilio code check if we're testing
	phonenumberCheckStatus := twilio.CheckPhoneNumber{
		Status: "approved",
	}
	if appconfig.Config.Gcloud.Project == "vama-prod" {
		status, twilioCheckErr := svc.twilio.Check(ctx, appconfig.Config.Twilio.Verify, url.Values{"To": []string{formattedNumber}, "Code": []string{req.Code}})
		if twilioCheckErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				400,
				http.StatusBadRequest,
				"Invalid verification code.",
				fmt.Sprintf("Invalid verification code. Err: %v", twilioCheckErr),
			)
		}
		phonenumberCheckStatus = *status
	}
	if phonenumberCheckStatus.Status != "approved" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Invalid verification code.",
			"Invalid verification code.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Could not begin transaction. Error: %v", txErr),
		)
	}

	// Defer a rollback unless we set commit to true
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	user, getUserByPhoneErr := svc.user.GetUserByPhone(ctx, tx, formattedNumber)
	if getUserByPhoneErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Error accessing db with phone number %s. Error: %v", formattedNumber, getUserByPhoneErr),
		)
	} else if user != nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"A user with this phone number already exists. Please try logging in.",
			"A user with this phone number already exists.",
		)
	}

	userName, usernameErr := svc.user.GetNextAvailableUsername(ctx, req.FirstName, req.LastName)
	if usernameErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Error generating unique username.  Err: %v", usernameErr),
		)
	} else if userName == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Error the randomly generated username returned nil for phone %s.", formattedNumber),
		)
	}
	userUpdateData := &response.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		CountryCode: req.CountryCode,
		Phonenumber: formattedNumber,
		Type:        "USER",
		Username:    *userName,
	}

	// Stripe protocol
	stripeParams := &stripe.CustomerParams{
		Phone: &userUpdateData.Phonenumber,
	}
	custmrID, createStripeCustomerErr := svc.CreateStripeCustomer(stripeParams)
	if createStripeCustomerErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Could not create Stripe user for phone number %s. Error: %v", userUpdateData.Phonenumber, createStripeCustomerErr),
		)
	}
	userUpdateData.StripeID = &custmrID

	// Save user in DB
	newUser, upsertUserErr := svc.user.UpsertUser(ctx, tx, userUpdateData)
	if upsertUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Could not create user for phone number %s. Error: %v", formattedNumber, upsertUserErr),
		)
	} else if newUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Error retrieving new user after upserting on signup with phone number %s.", formattedNumber),
		)
	}
	res, authenticateErr := svc.AuthUser(ctx, tx, newUser.UUID, newUser.ID)
	res.User = *newUser

	if authenticateErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Error retrieving access tokens for phone number %s. Error: %v", formattedNumber, authenticateErr),
		)
	}

	addUserToContactsErr := svc.AddUsersToContacts(ctx, tx, *newUser)
	if addUserToContactsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Error adding user's contacts: %v", addUserToContactsErr),
		)
	}

	initializePushSettingsErr := svc.push.InitializeSettings(ctx, newUser.Type, newUser.ID, tx)
	if initializePushSettingsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Error initializing push notification settings: %v", initializePushSettingsErr),
		)
	}

	// SendBird protocol
	sendbirdAccessToken, createSendbirdUserErr := svc.CreateSendBirdUser(*newUser)
	if createSendbirdUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("An error occurred getting sendbird user for phone number %s and UUID %s. Err: %v", formattedNumber, newUser.UUID, createSendbirdUserErr),
		)
	}

	defer func() {
		if !commit {
			signupErrMsg := fmt.Sprintf("A user w/ phone number %s failed to sign up and created a Sendbird user", formattedNumber)
			telegram.TelegramClient.SendMessage(signupErrMsg)
		}
	}()

	res.Credentials.SendBirdAccessToken = sendbirdAccessToken

	newSendBirdMetadata := request.SendBirdUserMetadata{
		Username: request.Username{Username: *userName},
		Type:     newUser.Type,
	}

	upsertMetadataErr := svc.sendbirdClient.UpsertUserMetadata(newUser.ID, newSendBirdMetadata)
	if upsertMetadataErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Error associating username with sendbird user %d. Error: %v", newUser.ID, upsertMetadataErr),
		)
	}

	// For deferred function to commit transaction instead of rollback
	commit = true
	welcomeErr := svc.vamaBot.SendWelcomeMessages(ctx, newUser.ID)
	if welcomeErr != nil {
		vlog.Errorf(ctx, "Error sending welcome messages: %v", welcomeErr)
	}
	return res, nil
}
