package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {

	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/search/global",
			Method:         "get",
			Handler:        SearchGlobal,
			RequestDecoder: decodeSearchGlobal,
		},
		{
			Path:           "/search/mention",
			Method:         "get",
			Handler:        SearchMention,
			RequestDecoder: decodeSearchMention,
		},
	}
	return res
}
