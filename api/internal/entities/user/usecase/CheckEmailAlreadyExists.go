package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (svc *usecase) CheckEmailAlreadyExists(ctx context.Context, runnable utils.Runnable, value string) (*response.User, error) {
	user, getUserByEmailErr := svc.repo.GetUserByEmail(ctx, runnable, value)
	if getUserByEmailErr != nil {
		return nil, getUserByEmailErr
	}

	return user, nil
}
