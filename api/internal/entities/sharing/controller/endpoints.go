package controller

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/sharing"
)

func GetBioData(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(sharing.Usecase)

	req := incomeRequest.(request.GetBioData)

	res, err := svc.GetBioData(ctx, req)

	return res, err
}

func GetLink(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(sharing.Usecase)

	req := incomeRequest.(request.GetLink)

	res, err := svc.GetLink(ctx, req)

	return res, err
}

func NewMessageLink(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(sharing.Usecase)

	req := incomeRequest.(request.NewMessageLink)

	res, err := svc.NewMessageLink(ctx, req)
	return res, err
}

func GetMessageByLinkSuffix(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(sharing.Usecase)

	req := incomeRequest.(request.GetMessageByLink)

	res, err := svc.GetMessageByLinkSuffix(ctx, req)
	return res, err
}

func PublicGetMessageByLink(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(sharing.Usecase)

	req := incomeRequest.(request.GetMessageByLink)

	res, err := svc.PublicGetMessageByLink(ctx, req)
	return res, err
}

func UpsertBioLinks(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(sharing.Usecase)

	req := incomeRequest.(request.UpsertBioLinks)

	res, err := svc.UpsertBioLinks(ctx, req)

	return res, err
}

func GetThemes(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(sharing.Usecase)

	req := incomeRequest.(request.GetThemes)

	res, err := svc.GetThemes(ctx, req)

	return res, err
}

func UpsertTheme(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(sharing.Usecase)

	req := incomeRequest.(request.UpsertTheme)

	res, err := svc.UpsertTheme(ctx, req)

	return res, err
}

func DeleteTheme(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(sharing.Usecase)

	req := incomeRequest.(request.DeleteTheme)

	err := svc.DeleteTheme(ctx, req)

	return struct{}{}, err
}
