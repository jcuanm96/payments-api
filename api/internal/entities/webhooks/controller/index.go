package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {

	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/stripe/event",
			Method:         "post",
			Handler:        StripeHandleEvent,
			RequestDecoder: decodeStripeHandleEvent,
		},
		{
			Path:           "/sendbird/event",
			Method:         "post",
			Handler:        SendbirdHandleEvent,
			RequestDecoder: decodeSendbirdHandleEvent,
		},
	}

	return res
}
