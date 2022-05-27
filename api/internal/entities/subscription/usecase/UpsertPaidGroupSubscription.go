package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (svc *usecase) UpsertPaidGroupSubscription(ctx context.Context, runnable utils.Runnable, subscription *response.PaidGroupChatSubscription) error {
	return svc.repo.UpsertPaidGroupSubscription(ctx, runnable, subscription)
}
