package service

import (
	"github.com/VamaSingapore/vama-api/internal/entities/search"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

type usecase struct {
	repo           search.Repository
	sendbirdClient sendbird.Client
	user           user.Usecase
}

func New(
	repo search.Repository,
	sendbirdClient sendbird.Client,
	user user.Usecase,
) search.Usecase {
	return &usecase{
		repo:           repo,
		sendbirdClient: sendbirdClient,
		user:           user,
	}
}
