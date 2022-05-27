package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errGettingChannel = "Something went wrong when trying to get channel."

func (svc *usecase) GetChannelWithUser(ctx context.Context, req request.GetChannelWithUser) (*sendbird.GroupChannel, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errGettingChannel,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errGettingChannel,
			"getCurrentUser returned nil for GetChannelWithUser",
		)
	}

	distinctMode := "distinct"
	showMember := true
	listGroupChannelsParams := sendbird.ListGroupChannelsParams{
		MembersExactlyIn: []string{strconv.Itoa(user.ID), strconv.Itoa(req.UserID)},
		DistinctMode:     &distinctMode,
		ShowMember:       &showMember,
	}
	channels, channelsErr := svc.sendbirdClient.ListGroupChannels(listGroupChannelsParams)
	if channelsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errGettingChannel,
			fmt.Sprintf("Error requesting group channels from sendbird: %v", channelsErr),
		)
	}

	var channel *sendbird.GroupChannel
	for _, c := range channels {
		if c.MemberCount == 2 {
			channel = &c
			break
		}
	}
	return channel, nil
}
