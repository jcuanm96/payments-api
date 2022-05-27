package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/webhooks"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const removeUserErr = "Something went wrong when trying to remove user from group."

func (svc *usecase) RemoveUserFromGroup(ctx context.Context, removedUserID int, channelID string) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			removeUserErr,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			removeUserErr,
			"user was nil in RemoveUserFromGroup",
		)
	}

	getGroupChannelParams := sendbird.GetGroupChannelParams{
		ShowMember: true,
	}
	channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(channelID, getGroupChannelParams)
	if getChannelErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			removeUserErr,
			fmt.Sprintf("Error getting Sendbird channel for removed user %d and channel %s in RemoveUserFromGroup: %v", removedUserID, channelID, getChannelErr),
		)
	} else if channel == nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			removeUserErr,
			fmt.Sprintf("Error the Sendbird channel %s came back as nil.", channelID),
		)
	}

	if !channel.HasOperator(fmt.Sprint(user.ID)) {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You cannot remove a user in this channel.",
			"You cannot remove a user in this channel.",
		)
	}

	isMember, isMemberErr := svc.sendbirdClient.IsMemberInGroupChannel(channel.ChannelURL, fmt.Sprint(user.ID))
	if isMemberErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			removeUserErr,
			fmt.Sprintf("Error something went wrong checking user %d membership in %s: %v", user.ID, channelID, isMemberErr),
		)
	}
	if !isMember.IsMember {
		return httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"This user does not exist in this group channel.",
			"This user does not exist in this group channel.",
		)
	}

	leaveParams := sendbird.LeaveGroupChannelParams{
		UserIDs: []string{fmt.Sprint(removedUserID)},
	}
	leaveErr := svc.sendbirdClient.LeaveGroupChannel(channelID, &leaveParams)
	if leaveErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			removeUserErr,
			fmt.Sprintf("Error something went wrong leaving Sendbird channel %s: %v", channelID, leaveErr),
		)
	}

	removedUser, getRemovedUserErr := svc.user.GetUserByID(ctx, removedUserID)
	if getRemovedUserErr != nil {
		vlog.Errorf(ctx, "Error getting the removed user %d in RemoveUserFromGroup with channel %s: %v", removedUser, channelID, getRemovedUserErr)
		removedUser = &response.User{}
	} else if removedUser == nil {
		vlog.Errorf(ctx, "The removed user %d in RemoveUserFromGroup returned nil in channel %s.", removedUserID, channelID)
		removedUser = &response.User{}
	}

	// Send message to channel
	message, data, rangesErr := webhooks.CalculateRemovedGroupAdminMessageRanges(user, removedUser)
	if rangesErr != nil {
		vlog.Errorf(ctx, "Error calculating ranges for remove group message: %v", rangesErr)
		return nil
	}
	sendMessageErr := svc.SendGroupMessage(ctx, channelID, user.ID, message, data, constants.ChannelEventRemoveCustomType)
	if sendMessageErr != nil {
		vlog.Errorf(ctx, "Error sending remove group message in %s: %v", channelID, sendMessageErr)
	}

	pushErr := svc.push.SendRemovedGroupNotification(ctx, removedUser, user, channel)
	if pushErr != nil {
		vlog.Errorf(ctx, "Error sending removed group push to %d: %v", removedUser.ID, pushErr)
	}

	return nil
}
