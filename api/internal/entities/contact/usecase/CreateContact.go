package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errAddingContact = "Something went wrong when adding a contact."

func (svc *usecase) CreateContact(ctx context.Context, contactID int) (*response.User, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errAddingContact,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errAddingContact,
		)
	}

	createContactErr := svc.repo.CreateContact(ctx, user.ID, contactID)
	if createContactErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errAddingContact,
			fmt.Sprintf("There was an issue creating a contact between users %d and user %d.  Err: %v", user.ID, contactID, createContactErr),
		)
	}
	res, userErr := svc.user.GetUserByID(ctx, contactID)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errAddingContact,
			fmt.Sprintf("There was an issue fetching user %d by ID.  Err: %v", contactID, userErr),
		)
	}
	if res == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errAddingContact,
			"The user you tried to add doesn't exist",
		)
	}
	return res, nil
}
