package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/webhooks"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) SendbirdCreateGroupEvent(ctx context.Context, event *sendbird.GroupChannelCreateEvent) error {
	if event.Channel.IsDistinct {
		return nil
	}

	// If backend makes the group, this ID is empty.
	// We assume the operator is the creator.
	if event.Inviter.UserID == "" {
		channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(event.Channel.ChannelURL, sendbird.GetGroupChannelParams{ShowMember: true})
		if getChannelErr != nil {
			vlog.Errorf(ctx, "Error getting channel %s: %s", event.Channel.ChannelURL, getChannelErr.Error())
			return getChannelErr
		}

		// If there is not exactly one operator in the channel,
		// we do not know who made the group.
		if len(channel.Operators) != 1 {
			return nil
		}

		event.Inviter = webhooks.ConvertUserToWebhookUser(&channel.Operators[0])
	}

	messagef := "%s created the group"
	message, data, rangesErr := webhooks.CalculateNicknameThenTextRanges(ctx, &event.Inviter, messagef)
	if rangesErr != nil {
		return rangesErr
	}

	return svc.chat.SendGroupMessage(ctx, event.Channel.ChannelURL, *data.Ranges[0].ID, message, data, constants.ChannelEventCustomType)
}
