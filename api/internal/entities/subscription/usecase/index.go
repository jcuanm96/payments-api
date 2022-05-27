package service

import (
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/stripe/stripe-go/v72/client"
)

type usecase struct {
	repo           subscription.Repository
	user           user.Usecase
	chat           *chat.Usecase
	stripeClient   *client.API
	sendbirdClient sendbird.Client
}

func New(
	repo subscription.Repository,
	user user.Usecase,
	chat *chat.Usecase,
	stripeClient *client.API,
	sendbirdClient sendbird.Client,
) subscription.Usecase {
	return &usecase{
		repo:           repo,
		user:           user,
		chat:           chat,
		stripeClient:   stripeClient,
		sendbirdClient: sendbirdClient,
	}
}
