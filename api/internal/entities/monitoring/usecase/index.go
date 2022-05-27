package service

import (
	monitoring "github.com/VamaSingapore/vama-api/internal/entities/monitoring"
)

type usecase struct {
	repo monitoring.Repository
}

func New(repo monitoring.Repository) monitoring.Usecase {
	return &usecase{
		repo: repo,
	}
}
