package controller

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
)

func UpdateFcmToken(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(push.Usecase)
	req := incomeRequest.(request.UpdateFcmToken)

	err := svc.UpdateUserFcmToken(ctx, req.Token)
	return nil, err
}

func SetGoatPostNotifications(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(push.Usecase)
	req := incomeRequest.(request.SetGoatPostNotifications)

	err := svc.SetGoatPostNotifications(ctx, req.GoatID, req.Enable)
	return nil, err
}

func SendTestNotification(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(push.Usecase)

	err := svc.SendTestNotification(ctx)
	return nil, err
}

func GetSettings(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(push.Usecase)

	res, err := svc.GetSettings(ctx)
	return res, err
}

func UpdateSetting(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(push.Usecase)
	req := incomeRequest.(request.UpdatePushSetting)

	err := svc.UpdateSetting(ctx, req)
	return nil, err
}
