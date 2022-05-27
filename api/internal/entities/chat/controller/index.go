package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {

	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/chat/goat/create",
			Method:         "post",
			Handler:        CreateGoatChatChannel,
			RequestDecoder: decodeCreateGoatChatChannel,
		},
		{
			Path:           "/chat/user/start",
			Method:         "post",
			Handler:        StartGoatChat,
			RequestDecoder: decodeStartGoatChat,
		},
		{
			Path:           "/chat/user/end",
			Method:         "post",
			Handler:        EndGoatChat,
			RequestDecoder: decodeEndGoatChat,
		},
		{
			Path:           "/chat/user/confirm",
			Method:         "post",
			Handler:        ConfirmGoatChat,
			RequestDecoder: decodeConfirmGoatChat,
		},
		{
			Path:           "/chat/messages/list",
			Method:         "get",
			Handler:        GetGoatChatPostMessages,
			RequestDecoder: decodeGetGoatChatPostMessages,
		},
		{
			Path:           "/chat/channel",
			Method:         "get",
			Handler:        GetChannelWithUser,
			RequestDecoder: decodeGetChannelWithUser,
		},
		{
			Path:           "/chat/paid/group/join",
			Method:         "post",
			Handler:        JoinPaidGroup,
			RequestDecoder: decodeJoinPaidGroup,
		},
		{
			Path:           "/chat/paid/group/create",
			Method:         "post",
			Handler:        CreatePaidGroupChannel,
			RequestDecoder: decodeCreatePaidGroupChannel,
		},
		{
			Path:           "/chat/paid/group/price",
			Method:         "patch",
			Handler:        UpdatePaidGroupPrice,
			RequestDecoder: decodeUpdatePaidGroupPrice,
		},
		{
			Path:           "/chat/paid/group/limits",
			Method:         "patch",
			Handler:        UpdatePaidGroupMemberLimits,
			RequestDecoder: decodeUpdatePaidGroupMemberLimits,
		},
		{
			Path:           "/chat/paid/group/metadata",
			Method:         "patch",
			Handler:        UpdatePaidGroupMetadata,
			RequestDecoder: decodeUpdateGroupMetadata,
		},
		{
			Path:           "/chat/paid/group/link",
			Method:         "patch",
			Handler:        UpdatePaidGroupLink,
			RequestDecoder: decodeUpdateGroupLink,
		},
		{
			Path:           "/chat/paid/group/leave",
			Method:         "post",
			Handler:        LeavePaidGroup,
			RequestDecoder: decodeLeaveGroup,
		},
		{
			Path:           "/chat/paid/group/cancel",
			Method:         "post",
			Handler:        CancelPaidGroup,
			RequestDecoder: decodeCancelPaidGroup,
		},
		{
			Path:           "/chat/paid/group/goat/list",
			Method:         "get",
			Handler:        ListGoatPaidGroups,
			RequestDecoder: decodeListGoatPaidGroups,
		},
		{
			Path:           "/chat/paid/group",
			Method:         "get",
			Handler:        GetPaidGroup,
			RequestDecoder: decodeGetPaidGroup,
		},
		{
			Path:           "/chat/paid/group",
			Method:         "delete",
			Handler:        DeletePaidGroup,
			RequestDecoder: decodeDeletePaidGroup,
		},
		{
			Path:           "/chat/paid/group/remove",
			Method:         "post",
			Handler:        ScheduledRemoveFromPaidGroup,
			RequestDecoder: decodeScheduledRemoveFromPaidGroup,
			Version:        constants.CLOUD_TASK_V1,
		},
		{
			Path:           "/chat/paid/group/ban",
			Method:         "post",
			Handler:        BanUserFromPaidGroup,
			RequestDecoder: decodeBanUserFromPaidGroup,
		},
		{
			Path:           "/chat/paid/group/ban",
			Method:         "delete",
			Handler:        UnbanUserFromPaidGroup,
			RequestDecoder: decodeUnbanUserFromPaidGroup,
		},
		{
			Path:           "/chat/ban/list",
			Method:         "get",
			Handler:        ListBannedUsers,
			RequestDecoder: decodeListBannedUsers,
		},
		{
			Path:           "/chat/link",
			Method:         "get",
			Handler:        GetDeepLinkInfo,
			RequestDecoder: decodeGetDeepLinkInfo,
		},
		{
			Path:           "/chat/link/check",
			Method:         "get",
			Handler:        CheckLinkSuffixIsTaken,
			RequestDecoder: decodeCheckLinkSuffixIsTaken,
		},
		{
			Path:           "/chat/group/member",
			Method:         "delete",
			Handler:        RemoveUserFromGroup,
			RequestDecoder: decodeRemoveUserFromGroup,
		},
		{
			Path:    "/chat/media/upload/url",
			Method:  "get",
			Handler: GetMediaURLs,
		},
		{
			Path:           "/chat/media/download/urls",
			Method:         "get",
			Handler:        GetBatchMediaDownloadURLs,
			RequestDecoder: decodeGetBatchMediaDownloadURLs,
		},
		{
			Path:           "/chat/media/verify",
			Method:         "get",
			Handler:        VerifyMediaObject,
			RequestDecoder: decodeVerifyMediaObject,
		},
		{
			Path:           "/chat/free/group/create",
			Method:         "post",
			Handler:        CreateFreeGroupChannel,
			RequestDecoder: decodeCreateFreeGroupChannel,
		},
		{
			Path:           "/chat/free/group/join",
			Method:         "post",
			Handler:        JoinFreeGroup,
			RequestDecoder: decodeJoinFreeGroup,
		},
		{
			Path:           "/chat/free/group",
			Method:         "get",
			Handler:        GetFreeGroup,
			RequestDecoder: decodeGetFreeGroup,
		},
		{
			Path:           "/chat/free/group/leave",
			Method:         "post",
			Handler:        LeaveFreeGroup,
			RequestDecoder: decodeLeaveGroup,
		},
		{
			Path:           "/chat/free/group/metadata",
			Method:         "patch",
			Handler:        UpdateFreeGroupMetadata,
			RequestDecoder: decodeUpdateGroupMetadata,
		},
		{
			Path:           "/chat/free/group/link",
			Method:         "patch",
			Handler:        UpdateFreeGroupLink,
			RequestDecoder: decodeUpdateGroupLink,
		},
		{
			Path:           "/chat/free/group/user/list",
			Method:         "get",
			Handler:        ListUserFreeGroups,
			RequestDecoder: decodeListUserFreeGroups,
		},
		{
			Path:           "/chat/free/group",
			Method:         "delete",
			Handler:        DeleteFreeGroup,
			RequestDecoder: decodeDeleteFreeGroup,
		},
		{
			Path:           "/chat/free/group/creators",
			Method:         "post",
			Handler:        AddFreeGroupCoCreators,
			RequestDecoder: decodeAddFreeGroupCoCreators,
		},
		{
			Path:           "/chat/free/group/creators",
			Method:         "delete",
			Handler:        RemoveFreeGroupCoCreator,
			RequestDecoder: decodeRemoveFreeGroupCoCreator,
		},
	}
	return res
}
