package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetUserByID(ctx context.Context, id int) (*response.User, error) {
	user, err := svc.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error searching for user %d: %v", id, err),
		)
	} else if user == nil {
		return nil, nil
	}
	return user, nil
}

func (svc *usecase) GetUserByEmail(ctx context.Context, runnable utils.Runnable, email string) (*response.User, error) {
	return svc.repo.GetUserByEmail(ctx, runnable, email)
}

func (svc *usecase) GetUserByPhone(ctx context.Context, runnable utils.Runnable, phone string) (*response.User, error) {
	return svc.repo.GetUserByPhone(ctx, runnable, phone)
}

func (svc *usecase) GetUserByStripeAccountID(ctx context.Context, runnable utils.Runnable, stripeAccountID string) (*response.User, error) {
	return svc.repo.GetUserByStripeAccountID(ctx, runnable, stripeAccountID)
}

func (svc *usecase) GetUserByStripeID(ctx context.Context, runnable utils.Runnable, stripeID string) (*response.User, error) {
	return svc.repo.GetUserByStripeID(ctx, runnable, stripeID)
}

func (svc *usecase) GetUserByUsername(ctx context.Context, runnable utils.Runnable, username string) (*response.User, error) {
	return svc.repo.GetUserByUsername(ctx, runnable, username)
}

func (svc *usecase) GetUserByUUID(ctx context.Context, runnable utils.Runnable, uuid string) (*response.User, error) {
	res, err := svc.repo.GetUserByUUID(ctx, runnable, uuid)
	if err != nil {
		return nil, err
	}
	return res, nil
}
