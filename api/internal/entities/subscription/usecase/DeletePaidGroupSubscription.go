package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (svc *usecase) DeletePaidGroupSubscription(ctx context.Context, runnable utils.Runnable, channelID string, userID int) error {
	return svc.repo.DeletePaidGroupSubscription(ctx, runnable, channelID, userID)
}
