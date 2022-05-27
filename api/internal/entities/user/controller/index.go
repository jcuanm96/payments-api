package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {
	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/users/me",
			Method:         "patch",
			Handler:        UpdateCurrentUser,
			RequestDecoder: decodeUpdateCurrentUser,
		},
		{

			Path:           "/users/me/phone",
			Method:         "patch",
			Handler:        UpdatePhone,
			RequestDecoder: decodeUpdatePhone,
		},
		{
			Path:    "/users/me",
			Method:  "get",
			Handler: GetCurrentUser,
		},
		{
			Path:    "/users/me",
			Method:  "get",
			Version: constants.API_V2,
			Handler: GetMyProfile,
		},
		{
			Path:           "/users",
			Method:         "get",
			Handler:        ListUsers,
			RequestDecoder: decodeListUsers,
		},
		{
			Path:           "/users/me/profile-avatar",
			Method:         "post",
			Handler:        UploadProfileAvatar,
			RequestDecoder: decodeUploadProfileAvatar,
		},
		{
			Path:           "/users/me/bio",
			Method:         "post",
			Handler:        UpsertBio,
			RequestDecoder: decodeUpsertBio,
		},
		{
			Path:           "/users/id",
			Method:         "get",
			Handler:        GetUserContactByID,
			RequestDecoder: decodeGetUserContactByID,
		},
		{
			Path:           "/users/goat/profile",
			Method:         "get",
			Handler:        GetGoatProfile,
			RequestDecoder: decodeGetGoatProfile,
		},
		{
			Path:           "/users/recs/goats",
			Method:         "get",
			Handler:        GetGoatUsers,
			RequestDecoder: decodeGetGoatUsers,
		},
		{
			Path:           "/users/block",
			Method:         "post",
			Handler:        BlockUser,
			RequestDecoder: decodeBlockUser,
		},
		{
			Path:           "/users/unblock",
			Method:         "delete",
			Handler:        UnblockUser,
			RequestDecoder: decodeUnblockUser,
		},
		{
			Path:           "/users/report",
			Method:         "post",
			Handler:        ReportUser,
			RequestDecoder: decodeReportUser,
		},
		{
			Path:           "/users/goat/upgrade",
			Method:         "patch",
			RequestDecoder: decodeUpgradeUserToGoat,
			Handler:        UpgradeUserToGoat,
		},
		{
			Path:    "/users/invitecodes",
			Method:  "get",
			Handler: GetInviteCodeStatuses,
		},
	}

	return res
}
