package controller

import (
	"context"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
)

func GetTransactions(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.GetTransactions)

	res, err := svc.GetTransactions(ctx, req)

	if err != nil {
		errMsg := fmt.Sprintf("GetTransactions returned an error. Err: %v", err)
		telegram.TelegramClient.SendMessage(errMsg)
	}

	return res, err
}

func GetPendingTransactions(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.GetTransactions)

	res, err := svc.GetPendingTransactions(ctx, req)

	if err != nil {
		errMsg := fmt.Sprintf("GetPendingTransactions returned an error. Err: %v", err)
		telegram.TelegramClient.SendMessage(errMsg)
	}

	return res, err
}

func MakeChatPaymentIntent(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.MakeChatPaymentIntent)

	err := svc.MakeChatPaymentIntent(ctx, req)

	return nil, err
}

func ConfirmPaymentIntent(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.ConfirmPaymentIntent)

	res, err := svc.ConfirmPaymentIntent(ctx, req)

	return res, err
}

func SaveDefaultPaymentMethod(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.DefaultPaymentMethod)

	res, err := svc.SaveDefaultPaymentMethod(ctx, req)

	return res, err
}

func GetMyPaymentMethods(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	res, err := svc.GetMyPaymentMethods(ctx)

	if err != nil {
		errMsg := fmt.Sprintf("GetMyPaymentMethods returned an error. Err: %v", err)
		telegram.TelegramClient.SendMessage(errMsg)
	}

	return res, err
}

func UpsertGoatChatsPrice(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.UpsertGoatChatsPrice)

	res, err := svc.UpsertGoatChatsPrice(ctx, req)

	return res, err
}

func UpsertBank(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.UpsertBank)

	err := svc.UpsertBank(ctx, req)

	return nil, err
}

func GetBalance(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.GetBalance)

	res, err := svc.GetBalance(ctx, req)

	return res, err
}

func GetGoatChatPrice(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.GetGoatChatPrice)

	res, err := svc.GetGoatChatPrice(ctx, req)

	if err != nil {
		errMsg := fmt.Sprintf("GetGoatChatPrice returned an error. Err: %v", err)
		telegram.TelegramClient.SendMessage(errMsg)
	}

	return res, err
}

func ListUnpaidProviders(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	res, err := svc.ListUnpaidProviders(ctx)
	return res, err
}

func MarkProviderAsPaid(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.MarkProviderAsPaid)

	err := svc.MarkProviderAsPaid(ctx, req)
	return nil, err
}

func GetPayoutPeriods(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.GetPayoutPeriods)

	res, err := svc.GetPayoutPeriods(ctx, req)

	return res, err
}

func ListPayoutHistory(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.ListPayoutHistory)

	res, err := svc.ListPayoutHistory(ctx, req)

	return res, err
}

func GetPayPeriod(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	req := incomeRequest.(request.GetPayPeriod)

	res, err := svc.GetPayPeriod(ctx, req)

	return res, err
}

func SendPendingBalanceNotifications(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(wallet.Usecase)

	err := svc.SendPendingBalanceNotifications(ctx)

	return nil, err
}
