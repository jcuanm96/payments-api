package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const ErrGettingMessageLink = "Something went wrong getting a link to share."

func (svc *usecase) NewMessageLink(ctx context.Context, req request.NewMessageLink) (*response.MessageLink, error) {
	var wg sync.WaitGroup

	var channel *sendbird.GroupChannel
	var getChannelErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		getChannelParams := sendbird.GetGroupChannelParams{}
		channel, getChannelErr = svc.sendbirdClient.GetGroupChannel(req.ChannelID, getChannelParams)
	}()

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			ErrGettingMessageLink,
			fmt.Sprintf("Error starting tx in NewMessageLink: %v", txErr),
		)
	}
	defer tx.Rollback(ctx)

	linkSuffix, generateLinkSuffixErr := svc.repo.GenerateMessageLinkSuffix(ctx, tx)
	if generateLinkSuffixErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			ErrGettingMessageLink,
			fmt.Sprintf("Error generating message link suffix: %v", generateLinkSuffixErr),
		)
	}

	insertMessageLinkErr := svc.repo.InsertMessageLink(ctx, tx, linkSuffix, req.MessageID, req.ChannelID)
	if insertMessageLinkErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			ErrGettingMessageLink,
			fmt.Sprintf("Error inserting new message link: %v", insertMessageLinkErr),
		)
	}

	wg.Wait()

	if getChannelErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			ErrGettingMessageLink,
			fmt.Sprintf("Error getting sendbird channel %s: %v", req.ChannelID, getChannelErr),
		)
	}

	// This will eventually be some map lookup for supported channel types
	if channel.CustomType != constants.ChannelTypePaidGroup {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You can't share messages from this chat.",
			fmt.Sprintf("Message sharing not supported for %s channels", channel.CustomType),
		)
	}

	redirectBaseURL := appconfig.Config.Gcloud.RedirectBaseURL
	res := &response.MessageLink{
		Link: fmt.Sprintf(constants.MESSAGE_BASE_URL_F, redirectBaseURL, linkSuffix),
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			ErrGettingMessageLink,
			fmt.Sprintf("Error committing tx: %v", commitErr),
		)
	}
	return res, nil
}
