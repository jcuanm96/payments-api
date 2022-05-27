package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) SetGoatPostNotifications(ctx context.Context, goatID int, enableNotifications bool) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when toggling notifications for creator's posts.",
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	} else if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"Something went wrong when toggling notifications for creator's posts.",
			"The user returned nil in SetGoatPostNotifications",
		)
	}

	token, getTokenErr := svc.repo.GetUserFcmToken(ctx, user.ID)
	if getTokenErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when toggling notifications for creator's posts.",
			fmt.Sprintf("Error getting fcm token from user %d: %v", user.ID, getTokenErr),
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when toggling notifications for creator's posts.",
			fmt.Sprintf("Could not begin transaction: %v", txErr),
		)
	}
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	if enableNotifications {
		enableNotificationsErr := svc.repo.EnableGoatPostNotifications(ctx, tx, user.ID, goatID)
		if enableNotificationsErr != nil {
			return httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when enabling notifications for creator's posts.",
				fmt.Sprintf("Error enabling notifications for creator %d for user %d: %v", goatID, user.ID, enableNotificationsErr),
			)
		}
		subscribeErr := svc.fcmClient.SubscribeToGoatPostTopic(ctx, goatID, token)
		if subscribeErr != nil {
			return httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when enabling notifications for creator's posts.",
				fmt.Sprintf("Error subscribing user %d to creator %d post topic: %v", user.ID, goatID, subscribeErr),
			)
		}
	} else {
		disableNotificationsErr := svc.repo.DisableGoatPostNotifications(ctx, tx, user.ID, goatID)
		if disableNotificationsErr != nil {
			return httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when disabling notifications for creator's posts.",
				fmt.Sprintf("Error disabling notifications for creator %d for user %d: %v", goatID, user.ID, disableNotificationsErr),
			)
		}
		unsubscribeErr := svc.fcmClient.UnsubscribeFromGoatPostTopic(ctx, goatID, token)
		if unsubscribeErr != nil {
			return httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when enabling notifications for creator's posts.",
				fmt.Sprintf("Error unsubscribing user %d from creator %d post topic: %v", user.ID, goatID, unsubscribeErr),
			)
		}
	}

	commit = true
	return nil
}
