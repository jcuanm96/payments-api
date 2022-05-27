package service

import (
	"github.com/VamaSingapore/vama-api/internal/entities/contact"
	"github.com/VamaSingapore/vama-api/internal/entities/search"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
)

type usecase struct {
	repo   contact.Repository
	user   user.Usecase
	search search.Usecase
}

func New(
	repo contact.Repository,
	user user.Usecase,
	search search.Usecase,
) contact.Usecase {
	return &usecase{
		repo:   repo,
		user:   user,
		search: search,
	}
}
