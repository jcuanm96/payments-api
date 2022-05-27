package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/search"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) GetAllSendbirdChannels(ctx context.Context, userID int, req request.SearchGlobal) (*search.Channels, error) {
	var wg sync.WaitGroup

	var globalChannels []sendbird.GroupChannel
	var getGlobalChannelsErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		globalChannels, getGlobalChannelsErr = svc.GetGlobalSendbirdChannels(userID, req)
	}()

	var myChannels []sendbird.GroupChannel
	var getMyChannelsErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		myChannels, getMyChannelsErr = svc.GetMyGroupSendbirdChannels(userID, req)
	}()

	var myDistinctChannels []sendbird.GroupChannel
	var getMyDistinctChannelsErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		myDistinctChannels, getMyDistinctChannelsErr = svc.GetMyDistinctSendbirdChannels(userID, req)
	}()

	var ownedByChannels []sendbird.GroupChannel
	var getOwnedByChannelsErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		ownedByChannels, getOwnedByChannelsErr = svc.GetSendbirdChannelsOwnedBy(ctx, userID, req)
	}()

	wg.Wait()

	if getGlobalChannelsErr != nil {
		return nil, getGlobalChannelsErr
	}

	if getMyChannelsErr != nil {
		return nil, getMyChannelsErr
	}

	if getMyDistinctChannelsErr != nil {
		return nil, getMyDistinctChannelsErr
	}

	if getOwnedByChannelsErr != nil {
		return nil, getOwnedByChannelsErr
	}

	myChannels = append(myChannels, myDistinctChannels...)

	channelsMap := map[string]struct{}{}
	svc.accumulateChannelMap(channelsMap, myChannels)
	svc.accumulateChannelMap(channelsMap, globalChannels)

	userIDStr := fmt.Sprint(userID)
	for _, channel := range ownedByChannels {
		// Remove duplicate channels
		_, channelExists := channelsMap[channel.ChannelURL]
		if len(globalChannels) > globalChannelsLimit && len(myChannels) > myChannelsLimit {
			break
		} else if !channelExists {
			isMember, isMemberErr := svc.sendbirdClient.IsMemberInGroupChannel(channel.ChannelURL, userIDStr)
			if isMemberErr != nil {
				vlog.Errorf(ctx, "Error checking isMember of channel %s in getAllSendbirdChannels: %v", channel.ChannelURL, isMemberErr)
				continue
			}
			if isMember.IsMember || channel.HasOperator(userIDStr) {
				myChannels = append(myChannels, channel)
			} else {
				globalChannels = append(globalChannels, channel)
			}
		}
	}

	if len(globalChannels) > globalChannelsLimit {
		globalChannels = globalChannels[:globalChannelsLimit]
	}

	if len(myChannels) > myChannelsLimit {
		myChannels = myChannels[:myChannelsLimit]
	}

	searchResults := search.Channels{
		GlobalChannels: globalChannels,
		MyChannels:     myChannels,
	}

	return &searchResults, nil
}

func (svc *usecase) accumulateChannelMap(acc map[string]struct{}, channels []sendbird.GroupChannel) {
	if acc == nil {
		return
	}

	for _, channel := range channels {
		acc[channel.ChannelURL] = struct{}{}
	}
}
