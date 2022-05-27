package service

import (
	"context"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/webhooks"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) SendbirdChangedGroupEvent(ctx context.Context, event *sendbird.GroupChannelChangedEvent) error {
	if event.Channel.IsDistinct {
		return nil
	}

	// event.Channel doesn't have many fields, such as operators.
	channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(event.Channel.ChannelURL, sendbird.GetGroupChannelParams{ShowMember: true})
	if getChannelErr != nil {
		vlog.Errorf(ctx, "Error getting channel %s: %v", event.Channel.ChannelURL, getChannelErr)
		return getChannelErr
	}

	// If there is not exactly one operator in the channel,
	// we do not know who made the change.
	if len(channel.Operators) != 1 {
		return nil
	}

	changer := webhooks.ConvertUserToWebhookUser(&channel.Operators[0])

	for _, change := range event.Changes {
		if change.Key == "name" {
			messagef := fmt.Sprintf("%%s changed the group name to %s", change.New)
			message, data, rangesErr := webhooks.CalculateNicknameThenTextRanges(ctx, &changer, messagef)
			if rangesErr != nil {
				continue
			}
			svc.chat.SendGroupMessage(ctx, event.Channel.ChannelURL, *data.Ranges[0].ID, message, data, constants.ChannelEventCustomType)
		} else if change.Key == "cover_url" {
			messagef := "%s updated the group photo"
			message, data, rangesErr := webhooks.CalculateNicknameThenTextRanges(ctx, &changer, messagef)
			if rangesErr != nil {
				return nil
			}
			svc.chat.SendGroupMessage(ctx, event.Channel.ChannelURL, *data.Ranges[0].ID, message, data, constants.ChannelEventCustomType)
		}
	}

	return nil
}
