package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errConfirmingChat = "Something went wrong when confirming creator chat."

func (svc *usecase) ConfirmGoatChat(ctx context.Context, req request.ConfirmGoatChat) (*response.ConfirmGoatChat, error) {
	var checkWG sync.WaitGroup
	var channel *chat.Channel
	var getChannelErr error
	var providerUser *response.User
	var userErr error

	checkWG.Add(1)
	go func() {
		defer checkWG.Done()
		providerUser, userErr = svc.user.GetCurrentUser(ctx)
	}()

	channel, getChannelErr = svc.GetChannel(ctx, req.SendBirdChannelID)
	if getChannelErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errConfirmingChat,
			fmt.Sprintf("Error getting channel members: %v", getChannelErr),
		)
	} else if channel == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errConfirmingChat,
			fmt.Sprintf("Channel %s came back as nil for ConfirmGoatChat", req.SendBirdChannelID),
		)
	}

	if channel.LastMessage == nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusInternalServerError,
			"The customer has not sent any messages",
			fmt.Sprintf("Sendbird channel %s is trying to be confirmed, but the customer hasn't sent any messages.", req.SendBirdChannelID),
		)
	}

	if channel.Data.ExpiresAt == nil {
		errMsg := fmt.Sprintf("No ExpiresAt in channel %s data", req.SendBirdChannelID)
		telegram.TelegramClient.SendMessage(errMsg)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errConfirmingChat,
			errMsg,
		)
	}

	if *channel.Data.ExpiresAt == constants.GOAT_CHAT_EXPIRES_AT_DEFAULT_VALUE {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"This chat has already ended. Please leave the chat and try again.",
			fmt.Sprintf("ExpiresAt was %d on confirm for channel %s", constants.GOAT_CHAT_EXPIRES_AT_DEFAULT_VALUE, req.SendBirdChannelID),
		)
	}

	if *channel.Data.ExpiresAt < time.Now().UnixMilli() {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Chat has expired, you have 7 days to respond to a chat request before it expires.",
			"Chat has expired, you have 7 days to respond to a chat request before it expires.",
		)
	}

	checkWG.Wait()

	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errConfirmingChat,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if providerUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errConfirmingChat,
			fmt.Sprintf("User was nil. Request: %v", req),
		)
	}

	if providerUser.Type != "GOAT" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can confirm a creator chat.",
			fmt.Sprintf("User %d not of type GOAT, was type %s", providerUser.ID, providerUser.Type),
		)
	}

	_, providerOk := channel.MemberUserIDs[providerUser.ID]
	_, customerOk := channel.MemberUserIDs[req.CustomerUserID]

	if !providerOk || !customerOk {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			errConfirmingChat,
			fmt.Sprintf("One or more users not in channel. providerOk: %v, customerOK: %v", providerOk, customerOk),
		)
	}

	startTS, getStartTSErr := svc.getStartMessageTS(ctx, req.SendBirdChannelID, req.CustomerUserID)
	if getStartTSErr == constants.ErrNotFound {
		return nil, httperr.NewCtx(
			ctx,
			constants.STATUS_FAILED_CONFIRM_CHAT,
			constants.STATUS_FAILED_CONFIRM_CHAT,
			errConfirmingChat,
			fmt.Sprintf("No StartTS found for %s", req.SendBirdChannelID),
		)
	} else if getStartTSErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errConfirmingChat,
			fmt.Sprintf("Error getting startTS for %s: %v", req.SendBirdChannelID, getStartTSErr),
		)
	}

	endGoatChatReq := request.EndGoatChat{
		SendBirdChannelID: req.SendBirdChannelID,
		StartTS:           startTS,
		IsPublic:          req.IsPublic,
	}

	endGoatChatRes, endGoatChatHttpErr := svc.EndGoatChat(ctx, endGoatChatReq)
	if endGoatChatHttpErr != nil {
		vlog.Errorf(ctx, "Error ending creator chat: %v", endGoatChatHttpErr)
		return nil, endGoatChatHttpErr
	}
	isStart := false
	state := constants.ChannelStateActive
	isPublic := endGoatChatRes.Post != nil
	updatedChannelData, channelDataErr := svc.UpdateChannelData(ctx, channel, isStart, state, isPublic)

	if channelDataErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errConfirmingChat,
			fmt.Sprintf("Error updating channel %s data: %v", req.SendBirdChannelID, channelDataErr),
		)
	}
	if updatedChannelData == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errConfirmingChat,
			fmt.Sprintf("Updated channel data came back nil for channel %s", req.SendBirdChannelID),
		)
	}

	confirmPaymentReq := request.ConfirmPaymentIntent{
		CustomerUserID: req.CustomerUserID,
	}

	_, confirmPaymentHttpErr := svc.wallet.ConfirmPaymentIntent(ctx, confirmPaymentReq)
	if confirmPaymentHttpErr != nil {
		vlog.Errorf(ctx, "Error confirming payment intent while confirming creator chat: %v", confirmPaymentHttpErr)
		return nil, confirmPaymentHttpErr
	}

	res := response.ConfirmGoatChat{
		ChannelData: *updatedChannelData,
	}

	if endGoatChatRes.Post != nil {
		res.Post = endGoatChatRes.Post
	}

	return &res, nil
}

func (svc *usecase) getStartMessageTS(ctx context.Context, channelID string, customerID int) (int64, error) {
	previousConversationEndTS, getEndTSErr := svc.repo.GetPreviousConversationEndTS(ctx, channelID)
	if getEndTSErr != nil {
		return 0, getEndTSErr
	}

	nextLimit := 1
	prevLimit := 0
	include := false
	customerIDStr := fmt.Sprint(customerID)
	listMessagesParams := sendbird.ListGroupChannelMessagesParams{
		ChannelURL: channelID,
		MessageTS:  previousConversationEndTS,
		NextLimit:  &nextLimit,
		PrevLimit:  &prevLimit,
		Include:    &include,
		SenderID:   &customerIDStr,
	}
	listMessagesResp, listMessagesErr := svc.sendbirdClient.ListGroupChannelMessages(listMessagesParams)
	if listMessagesErr != nil {
		return 0, listMessagesErr
	}

	if len(listMessagesResp.Messages) == 0 {
		return 0, constants.ErrNotFound
	}

	return listMessagesResp.Messages[0].CreatedAt, nil
}
