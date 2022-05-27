package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (svc *usecase) DeletePaidGroupProduct(ctx context.Context, runnable utils.Runnable, channelID string) error {
	return svc.repo.DeletePaidGroupProduct(ctx, runnable, channelID)
}
