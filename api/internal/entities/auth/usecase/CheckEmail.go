package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (svc *usecase) CheckEmail(ctx context.Context, req request.Email) (*response.Check, error) {
	res := response.Check{}

	user, checkEmailErr := svc.user.CheckEmailAlreadyExists(ctx, svc.repo.MasterNode(), req.Email)

	if checkEmailErr != nil {
		return nil, checkEmailErr
	}
	if user != nil && user.ID > 0 {
		res.UserID = user.ID
		res.IsTaken = true
		res.Message = "A user already exists with this email"
	}
	return &res, nil
}
