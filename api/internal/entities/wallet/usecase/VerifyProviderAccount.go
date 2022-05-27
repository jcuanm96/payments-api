package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) VerifyProviderAccount(ctx context.Context) (*response.VerifyProviderAccount, error) {
	currentUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error occurred when retrieving user from db. Err: %v", getCurrUserErr),
		)
	} else if currentUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			"Current user is nil for wallet.VerifyProviderAccount",
		)
	} else if currentUser.StripeAccountID == nil {
		res := response.VerifyProviderAccount{
			StripeAccountID:     "",
			HasDetailsSubmitted: false,
		}
		return &res, nil
	}

	return nil, nil
}
