package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	subscriptionrepo "github.com/VamaSingapore/vama-api/internal/entities/subscription/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) ScheduledRemoveFromPaidGroup(ctx context.Context, req cloudtasks.RemoveFromPaidGroupTask) error {
	subscription, getSubscriptionErr := subscriptionrepo.GetPaidGroupSubscription(ctx, svc.repo.MasterNode(), req.UserID, req.ChannelID)
	if getSubscriptionErr != nil {
		vlog.Errorf(ctx, "Error getting subscription in scheduled leave paid group user %d channel %s: %s", req.UserID, req.ChannelID, getSubscriptionErr.Error())
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to remove user from paid group chat.",
		)
	}

	// User must've resubscribed, don't remove
	if subscription != nil && time.Now().Before(subscription.CurrentPeriodEnd) {
		return nil
	}

	getGroupChannelParams := sendbird.GetGroupChannelParams{ShowMember: true}
	channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(req.ChannelID, getGroupChannelParams)

	// Channel no longer exists (may have been deleted), so user has already been removed
	if getChannelErr == sendbird.ErrGroupChannelNotFound {
		return nil
	} else if getChannelErr != nil {
		vlog.Errorf(ctx, "Error getting channel %s in ScheduledRemoveFromPaidGroup: %s", req.ChannelID, getChannelErr.Error())
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to remove user from paid group chat.",
		)
	}

	isMember, isMemberErr := svc.sendbirdClient.IsMemberInGroupChannel(channel.ChannelURL, fmt.Sprint(req.UserID))
	if isMemberErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to remove user from paid group chat.",
			fmt.Sprintf("Error something went wrong checking user %d membership in %s. Err: %v", req.UserID, channel.ChannelURL, isMemberErr),
		)
	}

	// User must've left the channel on their own, don't remove
	if !isMember.IsMember {
		return nil
	}

	leaveGroupChatParams := sendbird.LeaveGroupChannelParams{
		UserIDs: []string{fmt.Sprint(req.UserID)},
	}

	leaveGroupChatErr := svc.sendbirdClient.LeaveGroupChannel(req.ChannelID, &leaveGroupChatParams)
	if leaveGroupChatErr != nil {
		vlog.Errorf(ctx, "Error leaving Sendbird paid group chat for user %d and paid group channel %s: %s", req.UserID, req.ChannelID, leaveGroupChatErr.Error())
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to remove user from paid group chat.",
		)
	}

	return nil
}
