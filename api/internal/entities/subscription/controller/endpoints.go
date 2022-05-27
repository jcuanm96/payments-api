package controller

import (
	"context"
	"fmt"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
)

func SubscribeCurrUserToGoat(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(subscription.Usecase)
	req := incomeRequest.(request.SubscribeCurrUserToGoat)

	res, err := svc.SubscribeCurrUserToGoat(ctx, req)

	return res, err
}

func UnsubscribeCurrUserFromGoat(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(subscription.Usecase)
	req := incomeRequest.(request.UnsubscribeCurrUserFromGoat)

	res, err := svc.UnsubscribeCurrUserFromGoat(ctx, req)

	return res, err
}

func GetUserSubscriptions(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(subscription.Usecase)
	req := incomeRequest.(request.GetUserSubscriptions)

	res, err := svc.GetUserSubscriptions(ctx, req)

	return res, err
}

func GetSubscriptionPrices(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(subscription.Usecase)
	req := incomeRequest.(request.GetSubscriptionPrices)

	res, err := svc.GetSubscriptionPrices(ctx, req)

	return res, err
}

func UpdateSubscriptionPrice(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(subscription.Usecase)
	req := incomeRequest.(request.UpdateSubscriptionPrice)

	err := svc.UpsertGoatSubscriptionInfo(ctx, req.PriceInSmallestDenom, req.Currency)
	if err != nil {
		return nil, err
	}
	return &response.UpdateSubscriptionPrice{}, err
}

func CheckUserSubscribedToGoat(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(subscription.Usecase)

	req := incomeRequest.(request.CheckUserSubscribedToGoat)

	res, err := svc.CheckUserSubscribedToGoat(ctx, req)

	return res, err
}

func GetMyPaidGroupSubscriptions(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(subscription.Usecase)

	req := incomeRequest.(request.GetMyPaidGroupSubscriptions)

	res, err := svc.GetMyPaidGroupSubscriptions(ctx, req)

	return res, err
}

func GetMyPaidGroupSubscription(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(subscription.Usecase)

	req := incomeRequest.(request.GetMyPaidGroupSubscription)

	res, err := svc.GetMyPaidGroupSubscription(ctx, req)

	return res, err
}

func BatchUnsubscribePaidGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(subscription.Usecase)

	req := incomeRequest.(cloudtasks.StripePaidGroupUnsubscribeTask)

	batchUnsubscribeErr := svc.BatchUnsubscribePaidGroup(ctx, req)
	if batchUnsubscribeErr != nil {
		batchUnsubErrMsg := fmt.Sprintf("Error unsubscribing from paid group chat in BatchUnsubscribePaidGroup for channel %s. Batch IDs: %v. Err: %s", req.ChannelID, req.UserIDs, batchUnsubscribeErr.Error())
		telegram.TelegramClient.SendMessage(batchUnsubErrMsg)
	}

	return nil, batchUnsubscribeErr
}
