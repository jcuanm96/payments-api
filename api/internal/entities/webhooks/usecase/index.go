package service

import (
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/internal/entities/webhooks"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/stripe/stripe-go/v72/client"
)

type usecase struct {
	repo           webhooks.Repository
	stripeClient   *client.API
	sendbirdClient sendbird.Client
	wallet         wallet.Usecase
	subscription   subscription.Usecase
	chat           chat.Usecase
}

func New(
	repo webhooks.Repository,
	stripeClient *client.API,
	sendbirdClient sendbird.Client,
	wallet wallet.Usecase,
	subscription subscription.Usecase,
	chat chat.Usecase,

) webhooks.Usecase {
	return &usecase{
		repo:           repo,
		stripeClient:   stripeClient,
		sendbirdClient: sendbirdClient,
		wallet:         wallet,
		subscription:   subscription,
		chat:           chat,
	}
}
