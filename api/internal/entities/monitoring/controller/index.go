package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {

	res := []sconfiguration.EndpointConfiguration{
		{
			Path:            "/health/db",
			Method:          "get",
			Handler:         PingMainDatabase,
			ResponseEncoder: sconfiguration.SuccessResponse,
		},
		{
			Path:           "/version",
			Method:         "get",
			Handler:        GetMinRequiredVersionByDevice,
			RequestDecoder: decodeGetMinRequiredVersionByDevice,
		},
	}

	return res
}
