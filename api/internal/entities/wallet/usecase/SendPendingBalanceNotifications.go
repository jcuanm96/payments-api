package service

import (
	"context"
	"net/http"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const sendPendingBalancesDefaultErr = "Something went wrong sending pending balance notifications."

func (svc *usecase) SendPendingBalanceNotifications(ctx context.Context) error {
	pendingBalanceNotifications, getNotificationsErr := svc.repo.GetPendingBalanceNotifications(ctx, svc.repo.MasterNode())
	if getNotificationsErr != nil {
		vlog.Errorf(ctx, "Error getting the pending balances and fcm tokens from db: %v", getNotificationsErr)
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			sendPendingBalancesDefaultErr,
		)
	}

	svc.push.SendPendingBalanceNotifications(ctx, pendingBalanceNotifications)

	return nil
}
