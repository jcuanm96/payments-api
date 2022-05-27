package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errUnbanningUser = "Something went wrong when unbanning user from paid group."

func (svc *usecase) UnbanUserFromPaidGroup(ctx context.Context, bannedUserID int, channelID string) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"user was nil in UnbanUserFromPaidGroup",
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
			errUnbanningUser,
			fmt.Sprintf("Error getting Sendbird channel %s when unbanning user %d: %v", channelID, bannedUserID, getChannelErr),
		)
	} else if channel == nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUnbanningUser,
			fmt.Sprintf("Error the Sendbird channel %s came back as nil.", channelID),
		)
	}

	if !channel.HasOperator(fmt.Sprint(user.ID)) {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You cannot unban a user in this channel.",
			"You cannot unban a user in this channel.",
		)
	}

	unbanUserErr := svc.sendbirdClient.UnbanUserFromGroupChannel(channelID, fmt.Sprint(bannedUserID))
	if unbanUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUnbanningUser,
			fmt.Sprintf("Error unbanning user %d from channel %s in sendbird: %v", bannedUserID, channelID, unbanUserErr),
		)
	}

	deleteBanErr := svc.repo.DeleteBannedChatUser(ctx, bannedUserID, channelID, svc.repo.MasterNode())
	if deleteBanErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUnbanningUser,
			fmt.Sprintf("Error deleting banned user %d from channel %s by user %d: %v", bannedUserID, channelID, user.ID, deleteBanErr),
		)
	}

	return nil
}
