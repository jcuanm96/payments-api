package service

import (
	"github.com/VamaSingapore/vama-api/internal/entities/auth"
	"github.com/VamaSingapore/vama-api/internal/entities/follow"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/entities/search"
	user "github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/messaging"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/internal/upload"
	"github.com/kevinburke/twilio-go"

	"github.com/stripe/stripe-go/v72/client"
)

type usecase struct {
	repo           user.Repository
	push           *push.Usecase
	follow         *follow.Usecase
	sendBirdClient sendbird.Client
	msg            *messaging.Client
	gcsClient      *upload.Client
	twilio         *twilio.VerifyPhoneNumberService
	twilioVerifyID string
	stripeClient   *client.API
	auth           *auth.Usecase
	search         *search.Usecase
}

func New(
	repo user.Repository,
	push *push.Usecase,
	follow *follow.Usecase,
	sendbirdClient sendbird.Client,
	msg *messaging.Client,
	gcsClient *upload.Client,
	twilio *twilio.VerifyPhoneNumberService,
	twilioVerifyID string,
	stripeClient *client.API,
	auth *auth.Usecase,
	search *search.Usecase,
) user.Usecase {
	return &usecase{
		repo:           repo,
		push:           push,
		follow:         follow,
		sendBirdClient: sendbirdClient,
		msg:            msg,
		gcsClient:      gcsClient,
		twilio:         twilio,
		twilioVerifyID: twilioVerifyID,
		stripeClient:   stripeClient,
		auth:           auth,
		search:         search,
	}
}
