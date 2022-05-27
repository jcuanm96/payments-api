package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {

	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/userbio",
			Method:         "get",
			Handler:        GetBioData,
			RequestDecoder: decodeGetBioData,
			Version:        constants.PUBLIC_V1,
		},
		{
			Path:           "/userbio/link",
			Method:         "get",
			Handler:        GetLink,
			RequestDecoder: decodeGetLink,
			Version:        constants.PUBLIC_V1,
		},
		{
			Path:           "/sharing/message/link",
			Method:         "post",
			Handler:        NewMessageLink,
			RequestDecoder: decodeNewMessageLink,
		},
		{
			Path:           "/sharing/message/link",
			Method:         "get",
			Handler:        GetMessageByLinkSuffix,
			RequestDecoder: decodeGetMessageByLinkSuffix,
		},
		{
			Path:           "/sharing/message/link",
			Method:         "get",
			Handler:        PublicGetMessageByLink,
			RequestDecoder: decodePublicGetMessageByLink,
			Version:        constants.PUBLIC_V1,
		},
		{
			Path:           "/userbio/links",
			Method:         "patch",
			Handler:        UpsertBioLinks,
			RequestDecoder: decodeUpsertBioLinks,
		},
		{
			Path:           "/userbio/themes",
			Method:         "get",
			Handler:        GetThemes,
			RequestDecoder: decodeGetThemes,
		},
		{
			Path:           "/userbio/theme",
			Method:         "patch",
			Handler:        UpsertTheme,
			RequestDecoder: decodeUpsertTheme,
			Version:        constants.PUBLIC_V1,
		},
		/*{
			Path:           "/userbio/theme",
			Method:         "delete",
			Handler:        DeleteTheme,
			RequestDecoder: decodeDeleteTheme,
			Version:        constants.PUBLIC_V1,
		},*/
	}

	return res
}
