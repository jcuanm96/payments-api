package service

import "context"

func (svc *usecase) GetInviteCodeStatus(ctx context.Context, code string) (*int, bool, error) {
	return svc.repo.GetInviteCodeStatus(ctx, code)
}
