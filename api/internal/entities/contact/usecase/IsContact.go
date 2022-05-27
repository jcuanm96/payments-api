package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (svc *usecase) IsContact(ctx context.Context, contactID int) (*response.IsContact, error) {
	return svc.user.IsContact(ctx, contactID)
}
