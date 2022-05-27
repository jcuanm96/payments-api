package controller

import (
	"github.com/VamaSingapore/vama-api/internal/entities/auth"
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	"github.com/VamaSingapore/vama-api/internal/entities/contact"
	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	"github.com/VamaSingapore/vama-api/internal/entities/follow"
	"github.com/VamaSingapore/vama-api/internal/entities/monitoring"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/entities/search"
	"github.com/VamaSingapore/vama-api/internal/entities/sharing"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/internal/entities/webhooks"
	"github.com/VamaSingapore/vama-api/internal/logerr"
	"github.com/VamaSingapore/vama-api/internal/messaging"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/internal/token"
	"github.com/VamaSingapore/vama-api/internal/upload"
	stripe "github.com/stripe/stripe-go/v72/client"
)

type Ctr struct {
	tokenSvc     token.Service
	upload       upload.Uploader
	msg          messaging.Messenger
	Sendbird     sendbird.Client
	Lr           *logerr.Logerr
	Contact      contact.Usecase
	Monitoring   monitoring.Usecase
	User         user.Usecase
	Chat         chat.Usecase
	Auth         auth.Usecase
	Subscription subscription.Usecase
	Wallet       wallet.Usecase
	Feed         feed.Usecase
	Sharing      sharing.Usecase
	Webhooks     webhooks.Usecase
	Push         push.Usecase
	Follow       follow.Usecase
	Stripe       stripe.API
	Search       search.Usecase
}

func New(
	walletUsecase wallet.Usecase,
	monitoringUsecase monitoring.Usecase,
	contactUsecase contact.Usecase,
	authUsecase auth.Usecase,
	userUsecase user.Usecase,
	chatUsecase chat.Usecase,
	subscriptionUsecase subscription.Usecase,
	t token.Service,
	u upload.Uploader,
	msg messaging.Messenger,
	sendbird sendbird.Client,
	lr *logerr.Logerr,
	feedusecase feed.Usecase,
	sharingusecase sharing.Usecase,
	webhooksusecase webhooks.Usecase,
	pushusecase push.Usecase,
	followusecase follow.Usecase,
	stripe stripe.API,
	searchusecase search.Usecase,
) *Ctr {
	return &Ctr{
		Contact:      contactUsecase,
		Auth:         authUsecase,
		User:         userUsecase,
		Chat:         chatUsecase,
		Monitoring:   monitoringUsecase,
		Wallet:       walletUsecase,
		Subscription: subscriptionUsecase,
		tokenSvc:     t,
		upload:       u,
		msg:          msg,
		Sendbird:     sendbird,
		Lr:           lr,
		Feed:         feedusecase,
		Sharing:      sharingusecase,
		Webhooks:     webhooksusecase,
		Push:         pushusecase,
		Follow:       followusecase,
		Stripe:       stripe,
		Search:       searchusecase,
	}
}

func (ctr *Ctr) TokenSvc() token.Service {
	return ctr.tokenSvc
}
