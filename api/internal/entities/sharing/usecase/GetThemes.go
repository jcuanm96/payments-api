package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (svc *usecase) GetThemes(ctx context.Context, req request.GetThemes) (*response.GetThemes, error) {
	return svc.repo.GetThemes(ctx, req)
}
