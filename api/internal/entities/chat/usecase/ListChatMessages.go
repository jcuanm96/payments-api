package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) ListChatMessages(ctx context.Context, req request.ListChatMessages) (*sendbird.ListMessages, error) {
	channel, getChannelErr := svc.GetChannel(ctx, req.SendBirdChannelID)
	if getChannelErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Something went wrong when creating request to list users in Sendbird channel: %v", getChannelErr),
		)
	} else if channel == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			"Channel came back as nil for ListChatMessages",
		)
	}

	prevLimit := 0
	nextLimit := req.Limit
	listMessagesParams := sendbird.ListGroupChannelMessagesParams{
		ChannelURL: req.SendBirdChannelID,
		MessageTS:  req.MessageTsFrom,
		PrevLimit:  &prevLimit,
		NextLimit:  &nextLimit,
	}

	listMessages, listMessagesErr := svc.sendbirdClient.ListGroupChannelMessages(listMessagesParams)
	if listMessagesErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Something went wrong when listing Sendbird messages for channel %s: %v", req.SendBirdChannelID, listMessagesErr),
		)
	}
	if listMessages == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			"Sendbird response for listing messages was nil",
		)
	}

	if req.MessageTsTo == nil {
		return listMessages, nil
	}

	// Iterate from newest to oldest messages (largest to smallest timestamp)
	// The first message <= messageToTS is the last result
	// to be returned.
	for i := len(listMessages.Messages) - 1; i >= 0; i-- {
		message := listMessages.Messages[i]
		if message.CreatedAt <= *req.MessageTsTo {
			listMessages.Messages = listMessages.Messages[:i+1]
			return listMessages, nil
		}
	}

	return listMessages, nil
}
