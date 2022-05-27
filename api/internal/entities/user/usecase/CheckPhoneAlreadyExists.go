package service

import "context"

func (svc *usecase) CheckPhoneAlreadyExists(ctx context.Context, value string) (bool, error) {
	return svc.repo.CheckUserByFieldAlreadyExists(ctx, "phone_number", value)
}
