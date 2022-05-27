package service

import (
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

func (svc *usecase) GetGlobalSendbirdChannels(userID int, req request.SearchGlobal) ([]sendbird.GroupChannel, error) {
	publicMode := "public"
	requestLimit := int64(30) // Get a large number of channels in case there is overlap w/ my channels
	showMember := true
	listChannelsParams := sendbird.ListGroupChannelsParams{
		PublicMode:   &publicMode,
		NameContains: &req.Query,
		Limit:        &requestLimit,
		ShowMember:   &showMember,
	}

	channels, getChannelsErr := svc.sendbirdClient.ListGroupChannels(listChannelsParams)
	if getChannelsErr != nil {
		return nil, getChannelsErr
	}

	outputChannels := []sendbird.GroupChannel{}
	userIDStr := fmt.Sprint(userID)
	for _, channel := range channels {
		isMember, isMemberErr := svc.sendbirdClient.IsMemberInGroupChannel(channel.ChannelURL, userIDStr)
		if isMemberErr != nil {
			continue
		}
		if !isMember.IsMember {
			outputChannels = append(outputChannels, channel)
		}
	}

	return outputChannels, nil
}
