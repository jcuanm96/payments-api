package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errSendingTestNotification = "Something went wrong sending test notification."

func (svc *usecase) SendTestNotification(ctx context.Context) error {
	user, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not find user in the current context. Err: %v", getUserErr),
		)
	} else if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
			"user was nil in SendTestNotification",
		)
	}

	token, getTokenErr := svc.repo.GetUserFcmToken(ctx, user.ID)
	if getTokenErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errSendingTestNotification,
			fmt.Sprintf("Error getting fcm token from user %d: %v", user.ID, getTokenErr),
		)
	}

	title := "Test Backend Push Notification"
	body := "You're wonderful! :)"
	notificationErr := svc.fcmClient.PushNotification(ctx, token, title, body, map[string]string{})
	if notificationErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errSendingTestNotification,
			fmt.Sprintf("Error sending notification to user %d: %v", user.ID, notificationErr),
		)
	}

	return nil
}
