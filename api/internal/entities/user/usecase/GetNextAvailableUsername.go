package service

import (
	"context"
)

func (svc *usecase) GetNextAvailableUsername(ctx context.Context, firstName string, lastName string) (*string, error) {
	return svc.repo.GetNextAvailableUsername(ctx, firstName, lastName)
}
