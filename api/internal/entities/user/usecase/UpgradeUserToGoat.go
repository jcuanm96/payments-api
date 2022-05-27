package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (svc *usecase) UpgradeUserToGoat(ctx context.Context, formattedNumber string, req request.SignupSMSGoat) (*response.AuthSuccess, error) {
	return (*svc.auth).SignupSMSGoat(ctx, formattedNumber, req)
}
