package service

import (
	"context"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) GetSendbirdChannelsOwnedBy(ctx context.Context, userID int, req request.SearchGlobal) ([]sendbird.GroupChannel, error) {
	channelIDs, searchUsersDBErr := svc.repo.GetOwnedByChannels(ctx, req.Query, ownedByChannelsLimit, svc.repo.MasterNode())
	if searchUsersDBErr != nil {
		return nil, searchUsersDBErr
	}

	channels := []sendbird.GroupChannel{}
	var wg sync.WaitGroup
	wg.Add(len(channelIDs))
	for _, channelID := range channelIDs {
		go func(channelID string) {
			defer wg.Done()
			getParams := sendbird.GetGroupChannelParams{
				ShowMember: true,
			}
			channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(channelID, getParams)
			if getChannelErr != nil {
				vlog.Errorf(ctx, "Error getting Sendbird channel %s in get owned by search: %v", channelID, getChannelErr)
			} else if channel == nil {
				vlog.Errorf(ctx, "Sendbird channel %s returned nil in get owned by search.", channelID)
			} else {
				channels = append(channels, *channel)
			}
		}(channelID)
	}

	wg.Wait()

	return channels, nil
}
