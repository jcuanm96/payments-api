package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (svc *usecase) CheckUsername(ctx context.Context, req request.Username) (*response.Check, error) {
	res := response.Check{}

	user, checkUsernameErr := svc.user.CheckUsernameAlreadyExists(ctx, svc.repo.MasterNode(), req.Username)

	if checkUsernameErr != nil {
		return nil, checkUsernameErr
	}
	if user != nil {
		res.IsTaken = true
		res.Message = "A user already exists with this username"
	}
	return &res, nil
}
