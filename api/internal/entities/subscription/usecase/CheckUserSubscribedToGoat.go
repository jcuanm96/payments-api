package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) CheckUserSubscribedToGoat(ctx context.Context, req request.CheckUserSubscribedToGoat) (*response.CheckUserSubscribedToGoat, error) {
	user, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context: %v", getUserErr),
		)
	} else if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"user returned nil in CheckUserSubscribedToGoat",
		)
	}

	userSubscription, subscriptionErr := svc.GetUserSubscriptionByGoatID(ctx, user.ID, req.ProviderUserID)
	if subscriptionErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Failed when checking if user %d is already subscribed to creator %d: %v", user.ID, req.ProviderUserID, subscriptionErr),
		)
	} else if userSubscription == nil || time.Now().After(userSubscription.CurrentPeriodEnd) {
		return &response.CheckUserSubscribedToGoat{IsSubscribed: false}, nil
	}

	isSubscribed := userSubscription != nil

	res := response.CheckUserSubscribedToGoat{
		IsSubscribed:     isSubscribed,
		IsRenewing:       userSubscription.IsRenewing,
		CurrentPeriodEnd: userSubscription.CurrentPeriodEnd,
	}

	return &res, nil
}
