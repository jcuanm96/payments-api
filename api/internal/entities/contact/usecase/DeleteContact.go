package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const removeContactErr = "Something went wrong when removing contact."

func (svc *usecase) DeleteContact(ctx context.Context, contactID int) (*response.DeleteContact, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			removeContactErr,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			removeContactErr,
			"user returned nil in DeleteContact.",
		)
	}
	deleteErr := svc.repo.DeleteContact(ctx, user.ID, contactID)
	if deleteErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			removeContactErr,
			fmt.Sprintf("Could not delete contact between user %d and user %d. Err: %v", user.ID, contactID, deleteErr),
		)
	}
	return &response.DeleteContact{}, nil
}
