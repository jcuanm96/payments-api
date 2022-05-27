package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const removedNotificationErr = "Something went wrong sending removed user notification."

func (svc *usecase) SendRemovedGroupNotification(ctx context.Context, removedUser *response.User, removerUser *response.User, channel *sendbird.GroupChannel) error {
	token, getTokenErr := svc.repo.GetUserFcmToken(ctx, removedUser.ID)
	if getTokenErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			removedNotificationErr,
			fmt.Sprintf("Error getting fcm token from removed user %d: %v", removedUser.ID, getTokenErr),
		)
	}

	title := channel.Name
	removerUserNickname := fmt.Sprintf("%s %s", removerUser.FirstName, removerUser.LastName)
	body := fmt.Sprintf("%s removed you from the group", removerUserNickname)
	payload := map[string]string{
		"id":   channel.ChannelURL,
		"type": constants.ChannelEventRemoveCustomType,
	}
	notificationErr := svc.fcmClient.PushNotification(ctx, token, title, body, payload)
	if notificationErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			removedNotificationErr,
			fmt.Sprintf("Error sending notification to removed user %d: %v", removedUser.ID, notificationErr),
		)
	}

	return nil
}
