package controller

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/search"
)

func SearchGlobal(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(search.Usecase)

	req := incomeRequest.(request.SearchGlobal)

	res, err := svc.SearchGlobal(ctx, req)

	return res, err
}

func SearchMention(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(search.Usecase)

	req := incomeRequest.(request.SearchMention)

	res, err := svc.SearchMention(ctx, req)

	return res, err
}
