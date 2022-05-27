package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetSubscriptionPrices(ctx context.Context, req request.GetSubscriptionPrices) (*response.GoatSubscriptionInfo, error) {
	currUser, err := svc.user.GetCurrentUser(ctx)
	if err != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context: %v", err),
		)
	} else if currUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"currUser was nil in GetMyPaidGroupSubscriptions",
		)
	}

	res, getSubscriptionInfoErr := svc.repo.GetGoatSubscriptionInfo(ctx, req.GoatUserID)
	if getSubscriptionInfoErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Something went wrong getting creator %d's subscription price: %v", req.GoatUserID, getSubscriptionInfoErr),
		)
	}

	return res, nil
}
