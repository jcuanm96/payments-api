package service

import "context"

func (svc *usecase) SoftDeleteUser(ctx context.Context, userID int) error {
	return svc.repo.SoftDeleteUser(ctx, userID)
}
