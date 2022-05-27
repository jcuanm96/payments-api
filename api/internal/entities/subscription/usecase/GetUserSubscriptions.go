package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getUserSubscriptionsErr = "Something went wrong when trying to get your subscriptions."

func (svc *usecase) GetUserSubscriptions(ctx context.Context, req request.GetUserSubscriptions) (*response.GetUserSubscriptions, error) {
	currUser, err := svc.user.GetCurrentUser(ctx)
	if err != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getUserSubscriptionsErr,
			fmt.Sprintf("Could not find user in the current context: %v", err),
		)
	} else if currUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			getUserSubscriptionsErr,
			"currUser was nil in GetUserSubscriptions",
		)
	} else if currUser.StripeID == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getUserSubscriptionsErr,
			fmt.Sprintf("User %d does not have a Stripe customer ID", currUser.ID),
		)
	}

	subscriptions, lastID, repoErr := svc.repo.GetUserSubscriptions(ctx, currUser.ID, *req.CursorID, uint64(*req.Limit))
	if repoErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getUserSubscriptionsErr,
			fmt.Sprintf("Error getting provider list: %v", repoErr),
		)
	}

	res := response.GetUserSubscriptions{
		LastID:        *lastID,
		Subscriptions: subscriptions,
	}

	return &res, nil
}
