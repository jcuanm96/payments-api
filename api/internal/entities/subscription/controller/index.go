package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {

	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/subscriptions/prices",
			Method:         "get",
			RequestDecoder: decodeGetSubscriptionPrices,
			Handler:        GetSubscriptionPrices,
		},
		{
			Path:           "/subscriptions/feed/monthly",
			Method:         "post",
			Handler:        SubscribeCurrUserToGoat,
			RequestDecoder: decodeSubscribeCurrUserToGoat,
		},
		{
			Path:           "/subscriptions/me",
			Method:         "get",
			Handler:        GetUserSubscriptions,
			RequestDecoder: decodeGetUserSubscriptions,
		},
		{
			Path:           "/subscriptions/feed/monthly",
			Method:         "delete",
			Handler:        UnsubscribeCurrUserFromGoat,
			RequestDecoder: decodeUnsubscribeCurrUserFromGoat,
		},
		{
			Path:           "/subscriptions/prices/monthly/me",
			Method:         "patch",
			Handler:        UpdateSubscriptionPrice,
			RequestDecoder: decodeUpdateSubscriptionPrice,
		},
		{
			Path:           "/subscriptions/goats/check",
			Method:         "get",
			Handler:        CheckUserSubscribedToGoat,
			RequestDecoder: decodeCheckUserSubscribedToGoat,
		},
		{
			Path:           "/subscriptions/me/groups",
			Method:         "get",
			Handler:        GetMyPaidGroupSubscriptions,
			RequestDecoder: decodeGetMyPaidGroupSubscriptions,
		},
		{
			Path:           "/subscriptions/me/group",
			Method:         "get",
			Handler:        GetMyPaidGroupSubscription,
			RequestDecoder: decodeGetMyPaidGroupSubscription,
		},
		{
			Path:           "/subscriptions/groups/batch/unsubscribe",
			Method:         "post",
			Handler:        BatchUnsubscribePaidGroup,
			RequestDecoder: decodeBatchUnsubscribePaidGroup,
			Version:        constants.CLOUD_TASK_V1,
		},
	}

	return res
}
