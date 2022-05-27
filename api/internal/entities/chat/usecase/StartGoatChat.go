package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) StartGoatChat(ctx context.Context, req request.StartGoatChat) (*response.StartGoatChat, error) {
	var waitGroup sync.WaitGroup
	var user *response.User
	var userErr error

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		user, userErr = svc.user.GetCurrentUser(ctx)
	}()

	var channel *chat.Channel
	var sendBirdErr error

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		channel, sendBirdErr = svc.GetChannel(ctx, req.SendBirdChannelID)
	}()

	waitGroup.Wait()

	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to start a creator chat.  Please try again.",
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)

	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"Something went wrong when trying to start a creator chat.  Please try again.",
			"User was nil in StartGoatChat",
		)
	}

	if sendBirdErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to start a creator chat.  Please try again.",
			fmt.Sprintf("Something went wrong getting channel. Err: %v", sendBirdErr),
		)
	}

	if channel == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when starting creator chat",
			"Something went wrong getting channel. The channel returned nil for StartGoatChat",
		)
	}

	if _, ok := channel.MemberUserIDs[user.ID]; !ok {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You must be a member of the channel to start a chat with this creator.",
			fmt.Sprintf("User %d not in this channel", user.ID),
			"You must be a member of the channel to start a chat with this creator.",
		)
	}

	var providerUserID int
	for key := range channel.MemberUserIDs {
		if key != user.ID {
			providerUserID = key
			break
		}
	}

	tx, txErr := svc.repo.MasterNode().Begin(context.Background())
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to start a creator chat.  Please try again.",
			fmt.Sprintf("Error creating transaction: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	if channel.LastMessage != nil {
		updateGoatChatErr := svc.repo.UpdateGoatChatConversationEndTS(ctx, channel.LastMessage.CreatedAt, req.SendBirdChannelID, tx)
		if updateGoatChatErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Could not start creator chat",
				fmt.Sprintf("Error updating creator chat conversation_end_ts: %v", updateGoatChatErr),
			)
		}
	}

	insertGoatChatErr := svc.repo.InsertGoatChat(ctx, user.ID, providerUserID, req, tx)
	if insertGoatChatErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Could not start creator chat",
			fmt.Sprintf("Error inserting creator chat: %v", insertGoatChatErr),
		)
	}

	isStart := true
	state := constants.ChannelStatePending
	updatedChannelData, channelDataErr := svc.UpdateChannelData(ctx, channel, isStart, state, req.IsPublic)

	if channelDataErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to start a creator chat.  Please try again.",
			fmt.Sprintf("Error updating channel data: %v", channelDataErr),
		)
	}
	if updatedChannelData == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to start a creator chat.  Please try again.",
			fmt.Sprintf("Updated channel data came back nil: %v", req.SendBirdChannelID),
		)
	}

	defer func() {
		if !commit {
			isStart := false
			state := constants.ChannelStateActive
			_, channelDataErr := svc.UpdateChannelData(ctx, channel, isStart, state, req.IsPublic)
			if channelDataErr != nil {
				msg := fmt.Sprintf("Error resetting channel data for channel %s after failure in Start creator chat: %v", req.SendBirdChannelID, channelDataErr)
				vlog.Errorf(ctx, msg)
				telegram.TelegramClient.SendMessage(msg)
			}
		}
	}()

	paymentReq := request.MakeChatPaymentIntent{
		ProviderUserID: providerUserID,
	}

	paymentErr := svc.wallet.MakeChatPaymentIntent(ctx, paymentReq)

	if paymentErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Something went wrong when processing your payment to start a creator chat.  Please try again.",
			fmt.Sprintf("Error making payment while starting creator chat:  %v", paymentErr),
		)
	}

	commit = true
	res := &response.StartGoatChat{
		ChannelData: *updatedChannelData,
	}
	return res, nil
}
