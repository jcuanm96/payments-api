package repositories

import (
	"context"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (s *repository) GetUserByEmail(ctx context.Context, runnable utils.Runnable, email string) (*response.User, error) {
	return s.GetUser(ctx, runnable, "email = ?", email)
}

func (s *repository) GetUserByID(ctx context.Context, id int) (*response.User, error) {
	return s.GetUser(ctx, s.MasterNode(), "id = ?", id)
}

func (s *repository) GetUserByPhone(ctx context.Context, runnable utils.Runnable, phone string) (*response.User, error) {
	return s.GetUser(ctx, runnable, "phone_number = ?", phone)
}

func (s *repository) GetUserByStripeAccountID(ctx context.Context, runnable utils.Runnable, stripeAccountID string) (*response.User, error) {
	return s.GetUser(ctx, runnable, "stripe_account_id = ?", stripeAccountID)
}

func (s *repository) GetUserByStripeID(ctx context.Context, runnable utils.Runnable, stripeID string) (*response.User, error) {
	return s.GetUser(ctx, runnable, "stripe_id = ?", stripeID)
}

func (s *repository) GetUserByUsername(ctx context.Context, runnable utils.Runnable, username string) (*response.User, error) {
	return s.GetUser(ctx, runnable, "lower(username) = ?", strings.ToLower(username))
}

func (s *repository) GetUserByUUID(ctx context.Context, runnable utils.Runnable, uid string) (*response.User, error) {
	return s.GetUser(ctx, runnable, "uuid = ?", uid)
}
