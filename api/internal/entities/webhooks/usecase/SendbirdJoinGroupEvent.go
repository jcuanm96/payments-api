package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/webhooks"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) SendbirdJoinGroupEvent(ctx context.Context, event *sendbird.GroupChannelJoinEvent) error {
	if len(event.Users) == 0 || event.Channel.IsDistinct {
		return nil
	}

	// event.Channel doesn't have many fields, such as operators.
	channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(event.Channel.ChannelURL, sendbird.GetGroupChannelParams{ShowMember: true})
	if getChannelErr != nil {
		vlog.Errorf(ctx, "Error getting channel %s: %s", event.Channel.ChannelURL, getChannelErr.Error())
		return getChannelErr
	}

	joinedUsers := []sendbird.WebhookUser{}
	for _, joinedUser := range event.Users {
		// Likely a join event from group creation. Don't send any join messages
		if channel.HasOperator(joinedUser.UserID) {
			return nil
		}
		// If the user "invited" themselves, they just joined the channel.
		if joinedUser.Inviter.UserID == joinedUser.UserID {
			joinedUsers = append(joinedUsers, joinedUser)
		}
	}

	for _, joinedUser := range joinedUsers {
		messagef := "%s joined the group"
		message, data, rangesErr := webhooks.CalculateNicknameThenTextRanges(ctx, &joinedUser, messagef)
		if rangesErr != nil {
			continue
		}
		svc.chat.SendGroupMessage(ctx, event.Channel.ChannelURL, *data.Ranges[0].ID, message, data, constants.ChannelEventCustomType)
	}

	return nil
}
