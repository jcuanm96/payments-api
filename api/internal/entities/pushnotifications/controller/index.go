package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {
	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/push/fcm/token",
			Method:         "post",
			Handler:        UpdateFcmToken,
			RequestDecoder: decodeUpdateFcmToken,
		},
		{
			Path:           "/push/goat/posts",
			Method:         "post",
			Handler:        SetGoatPostNotifications,
			RequestDecoder: decodeSetGoatPostNotifications,
		},
		{
			Path:    "/push/test",
			Method:  "get",
			Handler: SendTestNotification,
		},
		{
			Path:    "/push/settings",
			Method:  "get",
			Handler: GetSettings,
		},
		{
			Path:           "/push/setting",
			Method:         "patch",
			Handler:        UpdateSetting,
			RequestDecoder: decodeUpdateSetting,
		},
	}

	return res
}
