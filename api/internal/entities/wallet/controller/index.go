package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {

	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/wallet/transactions/me",
			Method:         "get",
			Handler:        GetTransactions,
			RequestDecoder: decodeGetTransactions,
		},
		{
			Path:           "/wallet/transactions/pending",
			Method:         "get",
			Handler:        GetPendingTransactions,
			RequestDecoder: decodeGetPendingTransactions,
		},
		{
			Path:           "/wallet/payment-method/credit-card/me",
			Method:         "post",
			Handler:        SaveDefaultPaymentMethod,
			RequestDecoder: decodeDefaultPaymentMethod,
		},
		{
			Path:    "/wallet/payment-methods/me",
			Method:  "get",
			Handler: GetMyPaymentMethods,
		},
		{
			Path:           "/wallet/payment-intents/goat/chat",
			Method:         "post",
			Handler:        MakeChatPaymentIntent,
			RequestDecoder: decodeMakeChatPaymentIntent,
		},
		{
			Path:           "/wallet/payment-intents/goat/chat/confirm",
			Method:         "post",
			Handler:        ConfirmPaymentIntent,
			RequestDecoder: decodeConfirmPaymentIntent,
		},
		{
			Path:           "/wallet/goat/chat/me",
			Method:         "patch",
			Handler:        UpsertGoatChatsPrice,
			RequestDecoder: decodeUpsertGoatChatsPrice,
		},
		{
			Path:           "/wallet/payment-method/bank/me",
			Method:         "patch",
			Handler:        UpsertBank,
			RequestDecoder: decodeUpsertBank,
		},
		{
			Path:           "/wallet/balance/me",
			Method:         "get",
			Handler:        GetBalance,
			RequestDecoder: decodeGetBalance,
		},
		{
			Path:           "/wallet/goat/chat/price",
			Method:         "get",
			Handler:        GetGoatChatPrice,
			RequestDecoder: decodeGetGoatChatPrice,
		},
		{
			Path:           "/wallet/payout",
			Method:         "post",
			Version:        constants.PUBLIC_V1,
			Handler:        MarkProviderAsPaid,
			RequestDecoder: decodeMarkProviderAsPaid,
		},
		{
			Path:           "/wallet/payout/list",
			Method:         "get",
			Version:        constants.PUBLIC_V1,
			Handler:        ListUnpaidProviders,
			RequestDecoder: decodeListUnpaidProviders,
		},
		{
			Path:           "/wallet/payout/periods",
			Method:         "get",
			Version:        constants.PUBLIC_V1,
			Handler:        GetPayoutPeriods,
			RequestDecoder: decodeGetPayoutPeriods,
		},
		{
			Path:           "/wallet/payout/history",
			Method:         "get",
			Version:        constants.PUBLIC_V1,
			Handler:        ListPayoutHistory,
			RequestDecoder: decodeListPayoutHistory,
		},
		{
			Path:           "/wallet/payout/period",
			Method:         "get",
			Version:        constants.PUBLIC_V1,
			Handler:        GetPayPeriod,
			RequestDecoder: decodeGetPayPeriod,
		},
		{
			Path:    "/wallet/balance/notify",
			Method:  "post",
			Handler: SendPendingBalanceNotifications,
			Version: constants.CLOUD_SCHEDULER_V1,
		},
	}

	return res
}
