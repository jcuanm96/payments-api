package service

import (
	"context"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) SendPendingBalanceNotifications(ctx context.Context, notifications []wallet.PendingBalanceNotification) {
	title := "Vama"

	for _, notification := range notifications {
		body := fmt.Sprintf("You have a pending balance of $%.2f. Update your information to get paid!", float32(notification.AvailableBalance)/100.)
		payload := map[string]string{
			"type": constants.PendingBalanceNotificationType,
		}
		notificationErr := svc.fcmClient.PushNotification(ctx, notification.FcmToken, title, body, payload)
		if notificationErr != nil {
			vlog.Errorf(ctx, fmt.Sprintf("Error occurred sending pending balance notification: %v", notificationErr))
		}
	}
}
