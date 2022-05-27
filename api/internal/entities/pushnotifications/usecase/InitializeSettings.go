package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/jackc/pgx/v4"
)

func (svc *usecase) InitializeSettings(ctx context.Context, userType string, userID int, tx pgx.Tx) error {
	settings := push.UpdateSettings{}
	if userType == "GOAT" {
		settings.PendingBalance = constants.PUSH_NOTIFICATION_ON
	} else if userType == "USER" {
		settings.PendingBalance = constants.PUSH_NOTIFICATION_UNSET
	}
	return svc.repo.UpsertSettings(ctx, tx, userID, settings)
}
