package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) UpdateUserFcmToken(ctx context.Context, token string) error {
	currUser, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context: %v", getUserErr),
		)
	} else if currUser == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"currUser was nil in UpdateUserFcmToken",
		)
	}

	upsertErr := svc.repo.UpsertUserFcmToken(ctx, currUser.ID, token)
	if upsertErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error upserting user fcm token: %v", upsertErr),
		)
	}

	return nil
}
