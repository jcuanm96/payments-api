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

const gettingPaidGroupSubscriptionErr = "Something went wrong when trying to get your paid group chat subscription."

func (svc *usecase) GetMyPaidGroupSubscription(ctx context.Context, req request.GetMyPaidGroupSubscription) (*response.GetMyPaidGroupSubscription, error) {
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
			"user was nil in GetMyPaidGroupSubscription",
		)
	}

	subscription, getSubscriptionErr := svc.repo.GetPaidGroupSubscription(ctx, svc.repo.MasterNode(), user.ID, req.ChannelID)
	if getSubscriptionErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			gettingPaidGroupSubscriptionErr,
			fmt.Sprintf("Error getting paid group sub for user %d channel %s: %v", user.ID, req.ChannelID, getSubscriptionErr),
		)
	}

	if subscription == nil {
		return &response.GetMyPaidGroupSubscription{}, nil
	}

	paidGroup, getPaidGroupErr := (*svc.chat).GetPaidGroup(ctx, subscription.ChannelID)
	if getPaidGroupErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			gettingPaidGroupSubscriptionErr,
			fmt.Sprintf("Error getting paid group for user %d channel %s: %v", user.ID, subscription.ChannelID, getPaidGroupErr),
		)
	}
	subscription.Group = paidGroup.Group

	res := &response.GetMyPaidGroupSubscription{
		Subscription: subscription,
	}
	return res, nil
}
