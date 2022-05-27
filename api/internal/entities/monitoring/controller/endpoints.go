package controller

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/monitoring"
)

func PingMainDatabase(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(monitoring.Usecase)

	err := svc.PingMainDatabase(ctx)

	return nil, err
}

func GetMinRequiredVersionByDevice(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(monitoring.Usecase)

	req := incomeRequest.(request.GetMinRequiredVersionByDevice)

	res, err := svc.GetMinRequiredVersionByDevice(ctx, req)

	return res, err
}

func GetDashboard(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(monitoring.Usecase)

	res, err := svc.GetDashboard(ctx)

	return res, err
}
