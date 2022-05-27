package service

import (
	cloudTasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/internal/messaging"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/internal/upload"
	"github.com/stripe/stripe-go/v72/client"
)

type usecase struct {
	repo             chat.Repository
	feed             feed.Usecase
	user             user.Usecase
	wallet           wallet.Usecase
	subscription     subscription.Usecase
	msg              *messaging.Client
	sendbirdClient   sendbird.Client
	stripeClient     *client.API
	cloudTasksClient cloudTasks.Client
	gcsClient        *upload.Client
	push             push.Usecase
}

func New(
	repo chat.Repository,
	feed feed.Usecase,
	user user.Usecase,
	wallet wallet.Usecase,
	subscription subscription.Usecase,
	msg *messaging.Client,
	sendbirdClient sendbird.Client,
	stripeClient *client.API,
	cloudTasksClient cloudTasks.Client,
	gcsClient *upload.Client,
	push push.Usecase,
) chat.Usecase {
	return &usecase{
		repo:             repo,
		feed:             feed,
		user:             user,
		wallet:           wallet,
		subscription:     subscription,
		msg:              msg,
		sendbirdClient:   sendbirdClient,
		stripeClient:     stripeClient,
		cloudTasksClient: cloudTasksClient,
		gcsClient:        gcsClient,
		push:             push,
	}
}
