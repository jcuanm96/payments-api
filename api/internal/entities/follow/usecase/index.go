package service

import (
	"github.com/VamaSingapore/vama-api/internal/entities/follow"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
)

type usecase struct {
	repo follow.Repository
	user user.Usecase
}

func New(
	repo follow.Repository,
	user user.Usecase,
) follow.Usecase {
	return &usecase{
		repo: repo,
		user: user,
	}
}
