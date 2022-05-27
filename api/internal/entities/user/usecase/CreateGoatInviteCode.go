package service

import (
	"context"
)

func (svc *usecase) CreateGoatInviteCode(ctx context.Context, code string) (string, error) {
	err := svc.repo.CreateGoatInviteCode(ctx, code)
	if err != nil {
		return "", err
	}
	return code, nil
}
