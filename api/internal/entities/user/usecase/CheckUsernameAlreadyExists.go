package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (svc *usecase) CheckUsernameAlreadyExists(ctx context.Context, runnable utils.Runnable, username string) (*response.User, error) {
	user, getUsernameErr := svc.repo.GetUserByUsername(ctx, runnable, username)

	if getUsernameErr != nil {
		return nil, getUsernameErr
	}

	return user, nil
}
