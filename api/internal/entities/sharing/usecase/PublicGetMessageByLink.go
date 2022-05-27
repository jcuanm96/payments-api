package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/sharing"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/jackc/pgx/v4"
)

func (svc *usecase) PublicGetMessageByLink(ctx context.Context, req request.GetMessageByLink) (*response.PublicMessage, error) {
	deepLinkBaseURL := fmt.Sprintf("https://%s", appconfig.Config.Gcloud.RedirectBaseURL)
	deepLink := fmt.Sprintf(constants.MESSAGE_BASE_URL_F, deepLinkBaseURL, req.LinkSuffix)
	dynamicLink, getDynamicLinkErr := sharing.CreateDynamicLink(ctx, deepLink)
	if getDynamicLinkErr != nil {
		vlog.Errorf(ctx, "Error creating dynamic link for message suffix %s: %v", req.LinkSuffix, getDynamicLinkErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	}

	messageInfo, getMessageErr := svc.repo.GetMessageByLinkSuffix(ctx, svc.repo.MasterNode(), req.LinkSuffix)
	if getMessageErr == pgx.ErrNoRows {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"We couldn't find anything for this link.",
			"We couldn't find anything for this link.",
		)
	} else if getMessageErr != nil {
		vlog.Errorf(ctx, "Error getting message for suffix %s: %v", req.LinkSuffix, getMessageErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	}

	message, viewMessageErr := svc.sendbirdClient.ViewGroupChannelMessage(messageInfo.ChannelID, messageInfo.MessageID)
	if viewMessageErr != nil {
		vlog.Errorf(ctx, "Error getting message %s from channel %s in sendbird: %v", messageInfo.MessageID, messageInfo.ChannelID, viewMessageErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	}

	// This will change as we add support for all the message types
	if message.Type != sendbird.TextMessageType {
		return &response.PublicMessage{}, nil
	}

	senderID, atoiErr := strconv.Atoi(message.User.UserID)
	if atoiErr != nil {
		vlog.Errorf(ctx, "Error converting ID %s to int: %v", message.User.UserID, atoiErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	}
	sender, getSenderErr := svc.user.GetUserByID(ctx, senderID)
	if getSenderErr != nil {
		vlog.Errorf(ctx, "Error getting user %d by ID: %v", senderID, getSenderErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	}

	// Only keep the bare minimum fields
	returnedSenderUser := &response.User{
		Username:      sender.Username,
		ProfileAvatar: sender.ProfileAvatar,
		FirstName:     sender.FirstName,
		LastName:      sender.LastName,
	}

	res := &response.PublicMessage{
		DynamicLink: dynamicLink,
		Sender:      returnedSenderUser,
		MessageText: *message.Message,
	}
	return res, nil
}
