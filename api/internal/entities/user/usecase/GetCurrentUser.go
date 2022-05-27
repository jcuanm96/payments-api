package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/codes"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetCurrentUser(ctx context.Context) (*response.User, error) {
	userUUID := ""
	if ctx.Value("CURRENT_USER_UUID") != nil {
		userUUID = ctx.Value("CURRENT_USER_UUID").(string)
	}
	if userUUID == "" {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			"userUUID returned nil in GetCurrentUser",
		)
	}
	user, getUserErr := svc.repo.GetUserByUUID(ctx, svc.repo.MasterNode(), userUUID)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error getting user by uuid: %v", getUserErr),
		)
	}

	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"Something went wrong, or the user you're looking for doesn't exist. Please try again.",
			"user returned nil in GetCurrentUser",
		)
	}

	if user.DeletedAt != nil {
		return nil, httperr.NewCtx(
			ctx,
			codes.Omit,
			http.StatusForbidden,
			"User does not exist.",
			"User is deleted.",
		)
	}

	return user, nil
}
