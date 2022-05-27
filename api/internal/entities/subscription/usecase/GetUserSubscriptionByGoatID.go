package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (svc *usecase) GetUserSubscriptionByGoatID(ctx context.Context, userID int, goatUserID int) (*response.UserSubscription, error) {
	return svc.repo.GetUserSubscriptionByGoatID(ctx, userID, goatUserID)
}
