package service

import (
	"context"
)

func (svc *usecase) UpsertStripeAccountIDByEmail(ctx context.Context, email string, profileAvatarFilename string) error {
	return svc.repo.UpsertStripeAccountIDByEmail(ctx, email, profileAvatarFilename)
}