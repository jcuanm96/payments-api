package service

import (
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

func (svc *usecase) GetMyGroupSendbirdChannels(userID int, req request.SearchGlobal) ([]sendbird.GroupChannel, error) {
	distinctMode := "nondistinct"
	listChannelsParams := sendbird.ListMyGroupChannelsParams{
		NameContains: req.Query,
		Limit:        myChannelsLimit,
		DistinctMode: distinctMode,
	}

	channels, getChannelsErr := svc.sendbirdClient.ListMyGroupChannels(userID, listChannelsParams)
	if getChannelsErr != nil {
		return nil, getChannelsErr
	}

	return channels, nil
}
