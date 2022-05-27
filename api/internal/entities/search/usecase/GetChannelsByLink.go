package service

import (
	"context"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/search"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) GetChannelsByLink(ctx context.Context, userID int, req request.SearchGlobal) (*search.Channels, error) {
	req.Query = req.Query[1:]
	channelIDs, getChannelsErr := svc.repo.GetChannelsByLink(ctx, req.Query, channelLinkSearchLimit, svc.repo.MasterNode())
	if getChannelsErr != nil {
		return nil, getChannelsErr
	} else if len(channelIDs) <= 0 {
		results := search.Channels{
			GlobalChannels: []sendbird.GroupChannel{},
			MyChannels:     []sendbird.GroupChannel{},
		}
		return &results, nil
	}

	globalChannels := []sendbird.GroupChannel{}
	myChannels := []sendbird.GroupChannel{}
	userIDStr := fmt.Sprint(userID)
	for _, channelID := range channelIDs {
		getChannelParams := sendbird.GetGroupChannelParams{
			ShowMember: true,
		}
		channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(channelID, getChannelParams)
		if getChannelErr != nil {
			return nil, getChannelErr
		} else if channel == nil {
			vlog.Errorf(ctx, "Error tried searching for channel %s in db, but it doesn't exist in Sendbird", channelID)
			continue
		}

		if channel.HasOperator(userIDStr) {
			myChannels = append(myChannels, *channel)
			continue
		}

		isMember, isMemberErr := svc.sendbirdClient.IsMemberInGroupChannel(channel.ChannelURL, userIDStr)
		if isMemberErr != nil {
			continue
		}

		if isMember.IsMember {
			myChannels = append(myChannels, *channel)
		} else {
			globalChannels = append(globalChannels, *channel)
		}
	}

	results := search.Channels{
		GlobalChannels: globalChannels,
		MyChannels:     myChannels,
	}

	return &results, nil
}
