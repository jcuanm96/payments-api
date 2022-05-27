package service

import (
	cloudTasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	auth "github.com/VamaSingapore/vama-api/internal/entities/auth"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/internal/messaging"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/internal/token"
	"github.com/VamaSingapore/vama-api/internal/vamabot"
	"github.com/kevinburke/twilio-go"
	"github.com/stripe/stripe-go/v72/client"
)

const (
	ttlTokenHours        = 8
	ttlRefreshTokenHours = 24 * 365 * 100 //100 years
)

type usecase struct {
	repo             auth.Repository
	user             user.Usecase
	subscription     subscription.Usecase
	token            token.Service
	httpClient       *messaging.Client
	sendbirdClient   sendbird.Client
	vamaBot          vamabot.Client
	twilio           *twilio.VerifyPhoneNumberService
	stripeClient     *client.API
	cloudTasksClient cloudTasks.Client
	wallet           *wallet.Usecase
	push             push.Usecase
}

func New(
	repo auth.Repository,
	user user.Usecase,
	subscription subscription.Usecase,
	token token.Service,
	msg *messaging.Client,
	sendbirdClient sendbird.Client,
	vamaBot vamabot.Client,
	twilio *twilio.VerifyPhoneNumberService,
	stripeClient *client.API,
	cloudTasksClient cloudTasks.Client,
	wallet *wallet.Usecase,
	push push.Usecase,
) auth.Usecase {
	return &usecase{
		repo:             repo,
		user:             user,
		subscription:     subscription,
		token:            token,
		httpClient:       msg,
		sendbirdClient:   sendbirdClient,
		vamaBot:          vamaBot,
		twilio:           twilio,
		stripeClient:     stripeClient,
		cloudTasksClient: cloudTasksClient,
		wallet:           wallet,
		push:             push,
	}
}
