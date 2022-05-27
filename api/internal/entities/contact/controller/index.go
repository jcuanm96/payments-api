package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {

	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/contacts/me",
			Method:         "get",
			Handler:        GetContacts,
			RequestDecoder: decodeGetContacts,
		},
		{
			Path:           "/contacts/me/id",
			Method:         "get",
			Handler:        IsContact,
			RequestDecoder: decodeIsContact,
		},
		{
			Path:           "/contacts/me",
			Method:         "post",
			Handler:        CreateContact,
			RequestDecoder: decodeCreateContact,
		},
		{
			Path:           "/contacts/me",
			Method:         "delete",
			Handler:        DeleteContact,
			RequestDecoder: decodeDeleteContact,
		},
		{
			Path:           "/contacts/upload",
			Method:         "post",
			Handler:        UploadContacts,
			RequestDecoder: decodeUploadContacts,
		},
		{
			Path:    "/contacts/recommendations",
			Method:  "get",
			Handler: GetRecommendations,
		},
		{
			Path:           "/contacts/me/batch",
			Method:         "post",
			Handler:        BatchAddUsersToContacts,
			RequestDecoder: decodeBatchAddUsersToContacts,
			Version:        constants.CLOUD_TASK_V1,
		},
	}

	return res
}
