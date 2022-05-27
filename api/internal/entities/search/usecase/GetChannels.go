package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/search"
)

func (svc *usecase) GetChannels(ctx context.Context, userID int, req request.SearchGlobal) (*search.Channels, error) {
	// Search by link_suffix
	if string(req.Query[0]) == "@" {
		linkResults, getChannelsByLinkErr := svc.GetChannelsByLink(ctx, userID, req)
		if getChannelsByLinkErr != nil {
			return nil, getChannelsByLinkErr
		}
		return linkResults, nil
	}

	results, getSendbirdChannelsErr := svc.GetAllSendbirdChannels(ctx, userID, req)
	if getSendbirdChannelsErr != nil {
		return nil, getSendbirdChannelsErr
	}

	return results, nil
}
