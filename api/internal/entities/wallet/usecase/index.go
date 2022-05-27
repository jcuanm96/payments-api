package service

import (
	"github.com/VamaSingapore/vama-api/internal/entities/auth"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/stripe/stripe-go/v72/client"
)

type usecase struct {
	repo         wallet.Repository
	user         user.Usecase
	stripeClient *client.API
	auth         auth.Usecase
	push         push.Usecase
}

func New(repo wallet.Repository,
	user user.Usecase,
	stripeClient *client.API,
	auth auth.Usecase,
	push push.Usecase,
) wallet.Usecase {
	return &usecase{
		repo:         repo,
		user:         user,
		stripeClient: stripeClient,
		auth:         auth,
		push:         push,
	}
}
