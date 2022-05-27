package service

import (
	"context"
	"encoding/json"
	"fmt"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

func (svc *usecase) SendGroupMessage(ctx context.Context, channelURL string, fromUserID int, message string, data *response.AdminMessageData, customType string) error {
	dataBytes, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		vlog.Errorf(ctx, "Error marshaling group event message data: %v", marshalErr)
		return marshalErr
	}

	dataStr := string(dataBytes)
	sendMessageParams := &sendbird.SendMessageParams{
		MessageType: sendbird.TextMessageType,
		Message:     message,
		Data:        &dataStr,
		CustomType:  customType,
		UserID:      fmt.Sprint(fromUserID),
	}
	_, sendMessageErr := svc.sendbirdClient.SendMessage(channelURL, sendMessageParams)
	if sendMessageErr != nil {
		vlog.Errorf(ctx, "Error sending group event message: %v", sendMessageErr)
		return sendMessageErr
	}
	return nil
}
