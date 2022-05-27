package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/jackc/pgx/v4"
)

func (svc *usecase) UpsertUser(ctx context.Context, tx pgx.Tx, user *response.User) (*response.User, error) {
	newUser, upsertUserErr := svc.repo.UpsertUser(ctx, tx, user)
	if upsertUserErr != nil {
		return nil, upsertUserErr
	}
	return newUser, nil
}
