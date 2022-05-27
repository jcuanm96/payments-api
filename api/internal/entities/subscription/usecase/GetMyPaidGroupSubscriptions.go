package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getMyPaidGroupSubscriptionErr = "Something went wrong when trying to get your paid group chats."

func (svc *usecase) GetMyPaidGroupSubscriptions(ctx context.Context, req request.GetMyPaidGroupSubscriptions) (*response.GetMyPaidGroupSubscriptions, error) {
	user, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getMyPaidGroupSubscriptionErr,
			fmt.Sprintf("Could not find user in the current context: %v", getUserErr),
		)
	} else if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			getMyPaidGroupSubscriptionErr,
			"user was nil in GetMyPaidGroupSubscriptions",
		)
	}
	subscriptions, getSubscriptionsErr := svc.repo.GetMyPaidGroupSubscriptions(ctx, user.ID, req.CursorID, uint64(req.Limit))
	if getSubscriptionsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getMyPaidGroupSubscriptionErr,
			fmt.Sprintf("Error getting subscription list: %v", getSubscriptionsErr),
		)
	}

	var wg sync.WaitGroup
	wg.Add(len(subscriptions))
	for i := range subscriptions {
		go func(subscription *response.PaidGroupChatSubscription) {
			defer wg.Done()
			paidGroup, getPaidGroupErr := (*svc.chat).GetPaidGroup(ctx, subscription.ChannelID)
			if getPaidGroupErr != nil {
				vlog.Errorf(ctx, "Error getting paid group %s: %v", subscription.ChannelID, getPaidGroupErr)
				return
			}
			subscription.Group = paidGroup.Group
		}(&subscriptions[i])
	}
	wg.Wait()

	res := response.GetMyPaidGroupSubscriptions{
		Subscriptions: subscriptions,
	}

	return &res, nil
}
