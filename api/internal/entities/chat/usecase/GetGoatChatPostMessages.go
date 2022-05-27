package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetGoatChatPostMessages(ctx context.Context, req request.GetGoatChatPostMessages) (*response.GetGoatChatPostMessages, error) {
	conversation, getConversationErr := svc.repo.GetPostConversation(ctx, req.PostID)
	if getConversationErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error getting conversation from post with id %d: %v", req.PostID, getConversationErr),
		)
	}
	if conversation == nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not get conversation from post with id %d", req.PostID),
		)
	}

	if conversation.StartTS == 0 || conversation.ChannelID == "" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("The feed post w/ id %d is not a valid creator chat feed post.", req.PostID),
		)
	}

	fromTS := req.CursorMesssageTS
	if req.CursorMesssageTS == 0 {
		fromTS = int64(conversation.StartTS)
	}

	listChatMessagesReq := request.ListChatMessages{
		SendBirdChannelID: conversation.ChannelID,
		MessageTsFrom:     fromTS,
		Limit:             req.Limit,
	}

	if conversation.EndTS != 0 {
		endTS := int64(conversation.EndTS)
		listChatMessagesReq.MessageTsTo = &endTS
	}

	listChatMessagesRes, listMessagesHttpErr := svc.ListChatMessages(ctx, listChatMessagesReq)
	if listMessagesHttpErr != nil {
		vlog.Errorf(ctx, "Error getting Sendbird messages in GetGoatChatPostMessages: %v", listMessagesHttpErr)
		return nil, listMessagesHttpErr
	}

	res := response.GetGoatChatPostMessages{
		Messages: listChatMessagesRes.Messages,
	}

	return &res, nil
}
