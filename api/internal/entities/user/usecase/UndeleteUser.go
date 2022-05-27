package service

import "context"

func (svc *usecase) UndeleteUser(ctx context.Context, userID int) error {
	return svc.repo.UndeleteUser(ctx, userID)
}
