package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

func (svc *usecase) UpdateCurrentUser(ctx context.Context, req request.UpdateUser) (*response.User, error) {
	currentUser, getCurrentUserErr := svc.GetCurrentUser(ctx)
	if getCurrentUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error getting current user: %v", getCurrentUserErr),
		)
	} else if currentUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"user was nil in UpdateCurrentUser",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not begin transaction: %v", txErr),
		)
	}
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	if req.Email != nil {
		user, checkEmailErr := svc.CheckEmailAlreadyExists(ctx, tx, *req.Email)
		if checkEmailErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
				fmt.Sprintf("Error checking if email exists: %v", checkEmailErr),
			)
		}
		if user != nil && user.ID != currentUser.ID {
			return nil, httperr.NewCtx(
				ctx,
				400,
				http.StatusBadRequest,
				"The email you entered is unavailable.",
				"The email you entered is unavailable.",
			)
		}
	}

	if req.Username != nil {
		taken, isTakenErr := svc.repo.IsLinkSuffixTaken(ctx, tx, *req.Username)
		if isTakenErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
				fmt.Sprintf("Error checking if link is taken: %v", isTakenErr),
			)
		}
		// Taken and not changing your own username, like a case change: @David to @david
		if taken.IsTaken && !(taken.UserID != nil && *taken.UserID == currentUser.ID) {
			return nil, httperr.NewCtx(
				ctx,
				400,
				http.StatusBadRequest,
				"The username you entered is unavailable.",
				"The username you entered is unavailable.",
			)
		}
	}

	updateErr := svc.repo.UpdateUser(ctx, tx, currentUser.ID, req)
	if errors.Is(updateErr, constants.ErrNotFound) {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"user was not found in UpdateCurrentUser",
		)
	} else if updateErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error updating user %d: %v", currentUser.ID, updateErr),
		)
	}

	sendBirdParams := &sendbird.UpdateUserParams{}
	stripeParams := &stripe.CustomerParams{}

	newFirstName := currentUser.FirstName
	newLastName := currentUser.LastName

	if req.Email != nil {
		stripeParams.Email = req.Email
		currentUser.Email = req.Email
	}

	if req.FirstName != nil {
		newFirstName = *req.FirstName
		currentUser.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		newLastName = *req.LastName
		currentUser.LastName = *req.LastName
	}

	if req.ProfileAvatar != nil {
		sendBirdParams.ProfileURL = *req.ProfileAvatar
		currentUser.Phonenumber = *req.ProfileAvatar
	}

	if req.Username != nil {
		currentUser.Username = *req.Username
		sendBirdMetadataReq := request.SendBirdUserMetadata{
			Username: request.Username{Username: *req.Username},
			Type:     currentUser.Type,
		}
		sendBirdErr := svc.sendBirdClient.UpsertUserMetadata(currentUser.ID, sendBirdMetadataReq)
		if sendBirdErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
				fmt.Sprintf("Error updating sendbird metadata: %v", sendBirdErr),
			)
		}
	}

	// Update Stripe params
	if currentUser.StripeID != nil {
		stripeParams.Name = stripe.String(fmt.Sprintf("%s %s", newFirstName, newLastName))
		_, stripeUpdateErr := svc.stripeClient.Customers.Update(*currentUser.StripeID, stripeParams)

		if stripeUpdateErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
				fmt.Sprintf("Error updating Stripe customer for user: %d: %v", currentUser.ID, stripeUpdateErr),
			)
		}
	}

	// Update SendBird params
	sendBirdParams.Nickname = fmt.Sprintf("%s %s", newFirstName, newLastName)
	_, sendBirdErr := svc.sendBirdClient.UpdateUser(currentUser.ID, sendBirdParams)
	if sendBirdErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Something went wrong when updating SendBird user for id %d: %v", currentUser.ID, sendBirdErr),
		)
	}

	commit = true
	return currentUser, nil
}
