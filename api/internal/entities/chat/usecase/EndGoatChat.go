package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errEndingCreatorChat = "Something went wrong when ending creator chat."

func (svc *usecase) EndGoatChat(ctx context.Context, req request.EndGoatChat) (*response.EndGoatChat, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errEndingCreatorChat,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)

	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errEndingCreatorChat,
			fmt.Sprintf("User was nil. Request: %v", req),
		)
	}

	if user.Type != "GOAT" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can end a creator chat.",
			fmt.Sprintf("User %d not of type GOAT, was type %s", user.ID, user.Type),
		)
	}

	isMember, sendBirdErr := svc.IsUserInChannel(ctx, user.ID, req.SendBirdChannelID)

	if sendBirdErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errEndingCreatorChat,
			fmt.Sprintf("Error checking channel membership: %v", sendBirdErr),
		)
	}

	if !isMember {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			errEndingCreatorChat,
			fmt.Sprintf("User %d not in channel %s", user.ID, req.SendBirdChannelID),
		)
	}

	customerResult, repoErr := svc.repo.UpdateGoatChat(ctx, req)
	if repoErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errEndingCreatorChat,
			fmt.Sprintf("Error updating DB for end Goat chat: %v", repoErr),
		)
	}

	res := &response.EndGoatChat{}
	if !customerResult.IsPublic || !req.IsPublic {
		return res, nil
	}

	makeFeedPostReq := request.MakeFeedPost{
		GoatChatMessagesID: customerResult.GoatChatMessagesID,
	}

	feedPost, feedPostErr := svc.feed.MakeFeedPost(ctx, makeFeedPostReq)

	if feedPostErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong making a feed post after ending creator chat.",
			fmt.Sprintf("Error making feed post after ending creator chat: %v", feedPostErr),
		)
	}

	res.Post = feedPost

	return res, nil
}
