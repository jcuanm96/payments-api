package service

import (
	"context"

	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

type usecase struct {
	repo      push.Repository
	user      user.Usecase
	fcmClient push.Client
}

func New(ctx context.Context, repo push.Repository, user user.Usecase) push.Usecase {
	fcmClient, newFcmClientErr := push.NewClient(ctx)
	if newFcmClientErr != nil {
		vlog.Fatalf(ctx, "Error creating fcmClient: %v", newFcmClientErr)
	}
	return &usecase{
		repo:      repo,
		user:      user,
		fcmClient: fcmClient,
	}
}
