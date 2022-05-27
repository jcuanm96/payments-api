package chat

import (
	"context"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

type Usecase interface {
	CreateGoatChatChannel(ctx context.Context, req request.CreateGoatChatChannel) (*response.CreateGoatChatChannel, error)
	StartGoatChat(ctx context.Context, req request.StartGoatChat) (*response.StartGoatChat, error)
	EndGoatChat(ctx context.Context, req request.EndGoatChat) (*response.EndGoatChat, error)
	ConfirmGoatChat(ctx context.Context, req request.ConfirmGoatChat) (*response.ConfirmGoatChat, error)
	ListChatMessages(ctx context.Context, req request.ListChatMessages) (*sendbird.ListMessages, error)

	GetGoatChatPostMessages(ctx context.Context, req request.GetGoatChatPostMessages) (*response.GetGoatChatPostMessages, error)
	IsUserInChannel(ctx context.Context, userID int, sendBirdID string) (bool, error)
	GetChannel(ctx context.Context, sendBirdChannelID string) (*Channel, error)
	GetChannelWithUser(ctx context.Context, req request.GetChannelWithUser) (*sendbird.GroupChannel, error)
	UpdateChannelData(ctx context.Context, channel *Channel, resetData bool, state string, isPublic bool) (*response.ChannelData, error)

	JoinPaidGroup(ctx context.Context, req request.JoinPaidGroup) (*response.PaidGroupChatSubscription, error)
	GetPaidGroup(ctx context.Context, channelID string) (*response.GetPaidGroup, error)
	CreatePaidGroupChannel(ctx context.Context, req request.CreatePaidGroupChannel) (*response.PaidGroup, error)
	LeavePaidGroup(ctx context.Context, req request.LeaveGroup) (*response.LeavePaidGroup, error)
	CancelPaidGroup(ctx context.Context, req request.CancelPaidGroup) (*response.CancelPaidGroup, error)
	UpdatePaidGroupPrice(ctx context.Context, req request.UpdatePaidGroupPrice) error
	UpdatePaidGroupLink(ctx context.Context, req request.UpdateGroupLink) error
	UpdatePaidGroupMetadata(ctx context.Context, req request.UpdateGroupMetadata) error
	ListGoatPaidGroups(ctx context.Context, req request.ListGoatPaidGroups) (*response.ListGoatPaidGroups, error)
	DeletePaidGroup(ctx context.Context, req request.DeletePaidGroup) error
	ScheduledRemoveFromPaidGroup(ctx context.Context, req cloudtasks.RemoveFromPaidGroupTask) error
	BanUserFromPaidGroup(ctx context.Context, bannedUserID int, channelID string) error
	UnbanUserFromPaidGroup(ctx context.Context, bannedUserID int, channelID string) error
	RemoveUserFromGroup(ctx context.Context, removedUserID int, channelID string) error
	ListBannedUsers(ctx context.Context, req request.ListBannedUsers) (*response.ListBannedUsers, error)
	UpdatePaidGroupMemberLimits(ctx context.Context, req request.UpdatePaidGroupMemberLimits) error

	GetDeepLinkInfo(ctx context.Context, req request.GetDeepLinkInfo) (*response.DeepLinkInfo, error)
	CheckLinkSuffixIsTaken(ctx context.Context, req request.CheckLinkSuffixIsTaken) (*response.CheckLinkSuffixIsTaken, error)
	SendGroupMessage(ctx context.Context, channelURL string, fromUserID int, message string, data *response.AdminMessageData, customType string) error

	GetMediaURLs(ctx context.Context) (*response.MediaURLs, error)
	GetBatchMediaDownloadURLs(ctx context.Context, req request.GetBatchMediaDownloadURLs) (*response.BatchMediaDownloadURLs, error)
	VerifyMediaObject(ctx context.Context, req request.VerifyMediaObject) (*response.VerifyMediaObject, error)

	CreateFreeGroupChannel(ctx context.Context, req request.CreateFreeGroupChannel) (*response.FreeGroup, error)
	JoinFreeGroup(ctx context.Context, req request.JoinFreeGroup) error
	GetFreeGroup(ctx context.Context, channelID string) (*response.GetFreeGroup, error)
	LeaveFreeGroup(ctx context.Context, req request.LeaveGroup) error
	UpdateFreeGroupMetadata(ctx context.Context, req request.UpdateGroupMetadata) error
	UpdateFreeGroupLink(ctx context.Context, req request.UpdateGroupLink) error
	ListUserFreeGroups(ctx context.Context, req request.ListUserFreeGroups) (*response.ListUserFreeGroups, error)
	DeleteFreeGroup(ctx context.Context, req request.DeleteFreeGroup) error
	AddFreeGroupCoCreators(ctx context.Context, req request.AddFreeGroupChatCoCreators) error
	RemoveFreeGroupCoCreator(ctx context.Context, req request.RemoveFreeGroupChatCoCreator) error
}
