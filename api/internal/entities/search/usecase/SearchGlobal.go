package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const (
	globalChannelsLimit    = 5
	myChannelsLimit        = 5
	channelLinkSearchLimit = 5
	ownedByChannelsLimit   = 5
	userSearchLimit        = 5
	messageSearchLimit     = 10
)

const errSomethingWentWrongSearch = "Something went wrong when trying to do a global search."

func (svc *usecase) SearchGlobal(ctx context.Context, req request.SearchGlobal) (*response.SearchGlobal, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errSomethingWentWrongSearch,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errSomethingWentWrongSearch,
			"Current user was nil in SearchGlobal",
		)
	}

	var wg sync.WaitGroup

	var users []response.User
	var getUsersErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		users, getUsersErr = svc.GetUsers(ctx, req.Query, user.ID, userSearchLimit)
	}()

	var messages []response.SearchMessage
	var getMessagesErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		messages, getMessagesErr = svc.GetMessages(ctx, req.Query, user.ID, messageSearchLimit)
	}()

	channels, getChannelsErr := svc.GetChannels(ctx, user.ID, req)
	if getChannelsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errSomethingWentWrongSearch,
			fmt.Sprintf("Something went wrong when getting channels in global search for query %s: %v", req.Query, getChannelsErr),
		)
	} else if channels == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errSomethingWentWrongSearch,
			fmt.Sprintf("channels returned nil in global search for query %s.", req.Query),
		)
	}

	wg.Wait()

	if getUsersErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errSomethingWentWrongSearch,
			fmt.Sprintf("Something went wrong when getting users in global search for query %s: %v", req.Query, getUsersErr),
		)
	}

	if getMessagesErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errSomethingWentWrongSearch,
			fmt.Sprintf("Something went wrong when getting messages in global search for query %s: %v", req.Query, getMessagesErr),
		)
	}

	if len(channels.GlobalChannels) > globalChannelsLimit {
		channels.GlobalChannels = channels.GlobalChannels[:globalChannelsLimit]
	}

	res := response.SearchGlobal{
		GlobalChannels: channels.GlobalChannels,
		MyChannels:     channels.MyChannels,
		Users:          users,
		Messages:       messages,
	}

	return &res, nil
}
