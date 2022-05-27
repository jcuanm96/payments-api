package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

const signupCreatorErr = "Something went wrong when trying to sign up. Please try again."

func (svc *usecase) SignupSMSGoat(ctx context.Context, formattedNumber string, req request.SignupSMSGoat) (*response.AuthSuccess, error) {
	userID, wasFound, getInviteCodeStatusErr := svc.user.GetInviteCodeStatus(ctx, req.InviteCode)
	if getInviteCodeStatusErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Error retrieving creator invite code status for phone %s and code %s. Error: %v", formattedNumber, req.InviteCode, getInviteCodeStatusErr),
		)
	}

	if !wasFound || userID != nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"The invite code you entered is either invalid or taken.",
			"The invite code you entered is either invalid or taken.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Could not begin transaction. Error: %v", txErr),
		)
	}

	// Defer a rollback unless we set commit to true
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	currUser, getUserByPhoneErr := svc.user.GetUserByPhone(ctx, tx, formattedNumber)
	if getUserByPhoneErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Error accessing db for phone number %s. Error: %v", formattedNumber, getUserByPhoneErr),
		)
	}

	check, checkEmailErr := svc.CheckEmail(ctx, request.Email{Email: req.Email})
	if checkEmailErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Error occurred when checking if email exists. Err: %v", checkEmailErr),
		)
	} else if check == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			"check came back as nil when checking email in SignupSMSGoat",
		)
	} else if check.IsTaken && (currUser == nil || currUser.ID != check.UserID) { // Checking if user is upgrading and using their previous email
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Email %s has already been taken.  Please try a different email.", req.Email),
			fmt.Sprintf("Email %s has already been taken.  Please try a different email.", req.Email),
		)
	}

	username := strings.ReplaceAll(req.Username, " ", "")
	checkUsernameUser, checkUsernameExistsErr := svc.user.CheckUsernameAlreadyExists(ctx, tx, username)
	if checkUsernameExistsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Could not validate that signup username %s wasn't already taken. Err: %v", username, checkUsernameExistsErr),
		)
	} else if checkUsernameUser != nil && (currUser == nil || currUser.ID != checkUsernameUser.ID) { // Checking if user is upgrading and using their previous username
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			fmt.Sprintf("Username %s has already been taken. Please try a different username.", req.Username),
			fmt.Sprintf("Username %s has already been taken. Please try a different username.", req.Username),
		)
	}

	userUpdateData := &response.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		CountryCode: req.CountryCode,
		Phonenumber: formattedNumber,
		Username:    username,
		Email:       &req.Email,
		Type:        "GOAT",
	}

	var res *response.AuthSuccess

	// Initialize the creator user if they don't already have a user account
	// Else, upgrade their information to that of a creator
	if currUser == nil {
		// Stripe protocol (Create a Stripe customer ID)
		stripeParams := &stripe.CustomerParams{
			Email: userUpdateData.Email,
			Phone: &userUpdateData.Phonenumber,
		}
		customerID, createStripeCustomerErr := svc.CreateStripeCustomer(stripeParams)
		if createStripeCustomerErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				signupCreatorErr,
				fmt.Sprintf("Could not create Stripe user for phone number %s. Error: %v", formattedNumber, createStripeCustomerErr),
			)
		}
		userUpdateData.StripeID = &customerID

		// Save user in DB
		newUser, upsertUserErr := svc.user.UpsertUser(ctx, tx, userUpdateData)
		if upsertUserErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				signupCreatorErr,
				fmt.Sprintf("Could not create user for phone number %s. Error: %v", formattedNumber, upsertUserErr),
			)
		}

		if newUser == nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				signupCreatorErr,
				"newUser came back as nil in net new goat user signup",
			)
		}

		currUser = newUser

		// Authenticate user
		newAuthResponse, authenticateUserErr := svc.AuthUser(ctx, tx, currUser.UUID, currUser.ID)
		if authenticateUserErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				signupCreatorErr,
				fmt.Sprintf("Error retrieving access tokens for phone number %s. Error: %v", formattedNumber, authenticateUserErr),
			)
		}

		res = newAuthResponse

		// SendBird protocol (Create new SendBird user)
		sendbirdAccesToken, createSendBirdUserErr := svc.CreateSendBirdUser(*currUser)
		if createSendBirdUserErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				signupCreatorErr,
				fmt.Sprintf("An error occurred creating sendbird user for phone number %s and UUID %s. Error: %v", formattedNumber, currUser.UUID, createSendBirdUserErr),
			)
		}

		defer func() {
			if !commit {
				signupCreatorErrMsg := fmt.Sprintf("A net new user w/ phone number %s failed to sign up as a creator and created a Sendbird user", formattedNumber)
				telegram.TelegramClient.SendMessage(signupCreatorErrMsg)
			}
		}()

		res.Credentials.SendBirdAccessToken = sendbirdAccesToken
		res.User = *currUser
	} else {
		userUpdateData.ProfileAvatar = currUser.ProfileAvatar

		// Stripe protocol (Create a Stripe customer ID if the user doesn't already have one)
		if currUser.StripeID == nil {
			stripeParams := &stripe.CustomerParams{
				Email: stripe.String(*userUpdateData.Email),
				Phone: stripe.String(userUpdateData.Phonenumber),
			}
			customerID, createStripeCustomerErr := svc.CreateStripeCustomer(stripeParams)
			if createStripeCustomerErr != nil {
				return nil, httperr.NewCtx(
					ctx,
					500,
					http.StatusInternalServerError,
					signupCreatorErr,
					fmt.Sprintf("Could not create Stripe user for phone number %s. Error: %v", formattedNumber, createStripeCustomerErr),
				)
			}
			userUpdateData.StripeID = &customerID
		} else {
			// Update Stripe info
			updateStripeCustomerParams := &stripe.CustomerParams{
				Email: stripe.String(*userUpdateData.Email),
				Phone: stripe.String(userUpdateData.Phonenumber),
			}
			_, updateStripeCustomerErr := svc.stripeClient.Customers.Update(
				*currUser.StripeID,
				updateStripeCustomerParams,
			)

			if updateStripeCustomerErr != nil {
				return nil, httperr.NewCtx(
					ctx,
					500,
					http.StatusInternalServerError,
					signupCreatorErr,
					fmt.Sprintf("Could not update Stripe user for phone number %s. Error: %v", formattedNumber, updateStripeCustomerErr),
				)
			}

			userUpdateData.StripeID = currUser.StripeID
		}
		// Update user in DB
		updatedUser, upsertUserErr := svc.user.UpsertUser(ctx, tx, userUpdateData)
		if upsertUserErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				signupCreatorErr,
				fmt.Sprintf("Could not update user for phone number %s. Error: %v", formattedNumber, upsertUserErr),
			)
		} else if updatedUser == nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				signupCreatorErr,
				"updatedUser came back as nil when upgrading new goat user",
			)
		}

		// Authenticate user
		newAuthResponse, authenticateUserErr := svc.AuthUser(ctx, tx, currUser.UUID, currUser.ID)
		if authenticateUserErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				signupCreatorErr,
				fmt.Sprintf("Error retrieving access tokens for phone number %s. Error: %v", formattedNumber, authenticateUserErr),
			)
		}

		res = newAuthResponse
		res.User = *updatedUser

		sendbirdUser, getSendbirdUserErr := svc.sendbirdClient.GetUser(currUser.ID)
		if getSendbirdUserErr == nil {
			res.Credentials.SendBirdAccessToken = sendbirdUser.AccessToken
		} else {
			switch getSendbirdUserErr := getSendbirdUserErr.(type) {
			case sendbird.SendbirdErrorResponse:
				// If the user doesn't have a SendBird account, make one
				if getSendbirdUserErr.Code == sendbird.ResourceNotFound {
					sendbirdAccessToken, createSendBirdUserErr := svc.CreateSendBirdUser(*currUser)
					if createSendBirdUserErr != nil {
						return nil, httperr.NewCtx(
							ctx,
							500,
							http.StatusInternalServerError,
							signupCreatorErr,
							fmt.Sprintf("An error occurred creating sendbird user %d. Error: %v", currUser.ID, createSendBirdUserErr),
						)
					}
					res.Credentials.AccessToken = sendbirdAccessToken
				} else {
					return nil, httperr.NewCtx(
						ctx,
						500,
						http.StatusInternalServerError,
						signupCreatorErr,
						fmt.Sprintf("Error getting sendbird user ID %d: %v", currUser.ID, getSendbirdUserErr),
					)
				}
			default:
				return nil, httperr.NewCtx(
					ctx,
					500,
					http.StatusInternalServerError,
					signupCreatorErr,
					fmt.Sprintf("Error getting sendbird user ID %d: %v", currUser.ID, getSendbirdUserErr),
				)
			}
		}
	}

	defaultCurrency := constants.DEFAULT_CURRENCY
	defaultGoatChatPrice := int64(constants.MIN_GOAT_CHAT_PRICE_USD_IN_SMALLEST_DENOM)

	// Initialize a creator's chat price
	upsertGoatChatsPriceErr := svc.user.UpsertGoatChatsPrice(ctx, tx, defaultGoatChatPrice, defaultCurrency, currUser.ID)
	if upsertGoatChatsPriceErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Error initializing creator chat price for phone %s and code %s. Error: %v", formattedNumber, req.InviteCode, upsertGoatChatsPriceErr),
		)
	}

	// Credit creator for signing up
	payPeriod, getPayPeriodErr := (*svc.wallet).GetPayPeriod(ctx, request.GetPayPeriod{
		Timestamp: time.Now().Unix(),
	})
	if getPayPeriodErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Error getting pay period on creator signup for phone %s and code %s. Error: %v", formattedNumber, req.InviteCode, getPayPeriodErr),
		)
	}

	ledgerEntry := wallet.LedgerEntry{
		ProviderUserID:      currUser.ID,
		CustomerUserID:      constants.VAMA_USER_ID,
		StripeTransactionID: "",
		SourceType:          constants.TXN_HISTORY_SOURCE_TYPE_CHARGE,
		Amount:              constants.GOAT_SIGNUP_CREDIT_AMOUNT_USD_IN_SMALLEST_DENOM,
		VamaFee:             0,
		StripeFee:           0,
		CreatedTS:           0,
		BalanceDelta:        constants.GOAT_SIGNUP_CREDIT_AMOUNT_USD_IN_SMALLEST_DENOM,
		Currency:            constants.DEFAULT_CURRENCY,
		PayPeriodID:         payPeriod.ID,
	}
	insertTransactionErr := (*svc.wallet).InsertLedgerTransaction(ctx, tx, ledgerEntry)
	if insertTransactionErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Error inserting new transaction on creator signup for phone %s and code %s: %v", formattedNumber, req.InviteCode, insertTransactionErr),
		)
	}
	creditBalanceErr := (*svc.wallet).UpsertBalance(ctx, tx, ledgerEntry)
	if creditBalanceErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Error updating balance on creator signup for phone %s and code %s. Error: %v", formattedNumber, req.InviteCode, creditBalanceErr),
		)
	}

	initializePushSettingsErr := svc.push.InitializeSettings(ctx, currUser.Type, currUser.ID, tx)
	if initializePushSettingsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupErr,
			fmt.Sprintf("Error initializing push notification settings: %v", initializePushSettingsErr),
		)
	}

	// Use invite code
	useInviteCodeErr := svc.user.UseInviteCode(ctx, tx, req.InviteCode, currUser.ID)
	if useInviteCodeErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Error updating goat_invite_codes table for phone %s and code %s. Error: %v", formattedNumber, req.InviteCode, useInviteCodeErr),
		)
	}

	// SendBird protocol
	sendBirdMetadataReq := request.SendBirdUserMetadata{
		Username: request.Username{Username: username},
		Type:     "GOAT",
	}
	sendBirdMetadataErr := svc.sendbirdClient.UpsertUserMetadata(currUser.ID, sendBirdMetadataReq)
	if sendBirdMetadataErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Error associating username with sendbird user %d. Error: %v", currUser.ID, sendBirdMetadataErr),
		)
	}

	generateMyInviteCodesErr := svc.GenerateUserInviteCodes(ctx, tx, currUser.ID)
	if generateMyInviteCodesErr != nil {
		vlog.Errorf(ctx, "Error generating invite codes for user %d: %v", currUser.ID, generateMyInviteCodesErr)
	}

	addUserToContactsErr := svc.AddUsersToContacts(ctx, tx, *currUser)
	if addUserToContactsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			signupCreatorErr,
			fmt.Sprintf("Error adding user's contacts. Err: %v", addUserToContactsErr),
		)
	}

	commit = true
	welcomeErr := svc.vamaBot.SendWelcomeMessages(ctx, currUser.ID)
	if welcomeErr != nil {
		vlog.Errorf(ctx, "Error sending welcome messages: %v", welcomeErr)
	}
	return res, nil
}
