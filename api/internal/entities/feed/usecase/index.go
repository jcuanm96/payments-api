package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/upload"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

type usecase struct {
	repo         feed.Repository
	user         user.Usecase
	push         *push.Usecase
	chat         *chat.Usecase
	subscription subscription.Usecase
	*upload.Client
	gcsClient *upload.Client
}

func New(
	repo feed.Repository,
	user user.Usecase,
	push *push.Usecase,
	chat *chat.Usecase,
	subscription subscription.Usecase,
	gcsClient *upload.Client,
) feed.Usecase {
	if chat == nil {
		vlog.Fatalf(context.Background(), "chat usecase was nil in feed constructor")
	}
	return &usecase{
		repo:         repo,
		user:         user,
		push:         push,
		chat:         chat,
		subscription: subscription,
		gcsClient:    gcsClient,
	}
}
