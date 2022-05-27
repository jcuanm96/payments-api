package service

import (
	"github.com/VamaSingapore/vama-api/internal/entities/sharing"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/internal/upload"
)

type usecase struct {
	repo           sharing.Repository
	user           user.Usecase
	sendbirdClient sendbird.Client
	gcsClient      *upload.Client
}

func New(
	repo sharing.Repository,
	user user.Usecase,
	sendbirdClient sendbird.Client,
	gcsClient *upload.Client,
) sharing.Usecase {
	return &usecase{
		repo:           repo,
		user:           user,
		sendbirdClient: sendbirdClient,
		gcsClient:      gcsClient,
	}
}
