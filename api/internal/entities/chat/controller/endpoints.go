package controller

import (
	"context"
	"fmt"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
)

func CreateGoatChatChannel(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.CreateGoatChatChannel)

	res, err := svc.CreateGoatChatChannel(ctx, req)

	return res, err
}

func StartGoatChat(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.StartGoatChat)

	res, err := svc.StartGoatChat(ctx, req)

	return res, err
}

func EndGoatChat(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.EndGoatChat)

	res, err := svc.EndGoatChat(ctx, req)

	return res, err
}

func ConfirmGoatChat(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.ConfirmGoatChat)

	res, err := svc.ConfirmGoatChat(ctx, req)

	return res, err
}

func GetGoatChatPostMessages(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.GetGoatChatPostMessages)

	res, err := svc.GetGoatChatPostMessages(ctx, req)

	return res, err
}

func GetChannelWithUser(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.GetChannelWithUser)

	channel, err := svc.GetChannelWithUser(ctx, req)
	res := &response.GetChannelWithUser{}
	if channel != nil {
		res.ChannelID = channel.ChannelURL
	}

	return res, err
}

func CreatePaidGroupChannel(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.CreatePaidGroupChannel)

	res, err := svc.CreatePaidGroupChannel(ctx, req)

	return res, err
}

func JoinPaidGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.JoinPaidGroup)

	res, err := svc.JoinPaidGroup(ctx, req)
	return res, err
}

func UpdatePaidGroupPrice(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.UpdatePaidGroupPrice)

	err := svc.UpdatePaidGroupPrice(ctx, req)

	return nil, err
}

func UpdatePaidGroupMetadata(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.UpdateGroupMetadata)

	err := svc.UpdatePaidGroupMetadata(ctx, req)

	return nil, err
}

func UpdatePaidGroupLink(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.UpdateGroupLink)

	err := svc.UpdatePaidGroupLink(ctx, req)

	return nil, err
}

func LeavePaidGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.LeaveGroup)

	res, err := svc.LeavePaidGroup(ctx, req)

	return res, err
}

func CancelPaidGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.CancelPaidGroup)

	res, err := svc.CancelPaidGroup(ctx, req)

	return res, err
}

func ScheduledRemoveFromPaidGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(cloudtasks.RemoveFromPaidGroupTask)

	scheduledErr := svc.ScheduledRemoveFromPaidGroup(ctx, req)
	if scheduledErr != nil {
		scheduledErrMsg := fmt.Sprintf("Error occurred when processing task from ScheduledRemoveFromPaidGroup for channel %s. user id: %d. Err: %s", req.ChannelID, req.UserID, scheduledErr.Error())
		telegram.TelegramClient.SendMessage(scheduledErrMsg)
	}

	return nil, scheduledErr
}

func ListGoatPaidGroups(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.ListGoatPaidGroups)

	res, err := svc.ListGoatPaidGroups(ctx, req)
	return res, err
}

func DeletePaidGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.DeletePaidGroup)

	err := svc.DeletePaidGroup(ctx, req)
	return nil, err
}

func GetPaidGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.GetPaidGroup)

	res, err := svc.GetPaidGroup(ctx, req.ChannelID)
	return res, err
}

func BanUserFromPaidGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.BanUserFromPaidGroup)

	err := svc.BanUserFromPaidGroup(ctx, req.BannedUserID, req.ChannelID)
	return nil, err
}

func UnbanUserFromPaidGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.UnbanUserFromPaidGroup)

	err := svc.UnbanUserFromPaidGroup(ctx, req.BannedUserID, req.ChannelID)
	return nil, err
}

func CheckLinkSuffixIsTaken(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.CheckLinkSuffixIsTaken)

	res, err := svc.CheckLinkSuffixIsTaken(ctx, req)
	return res, err
}

func RemoveUserFromGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.RemoveUserFromGroup)

	err := svc.RemoveUserFromGroup(ctx, req.UserID, req.ChannelID)
	return nil, err
}

func ListBannedUsers(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.ListBannedUsers)

	res, err := svc.ListBannedUsers(ctx, req)
	return res, err
}

func UpdatePaidGroupMemberLimits(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)
	req := incomeRequest.(request.UpdatePaidGroupMemberLimits)

	err := svc.UpdatePaidGroupMemberLimits(ctx, req)

	return nil, err
}

func GetDeepLinkInfo(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)
	req := incomeRequest.(request.GetDeepLinkInfo)

	res, err := svc.GetDeepLinkInfo(ctx, req)
	return res, err
}

func GetMediaURLs(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	res, err := svc.GetMediaURLs(ctx)
	return res, err
}

func GetBatchMediaDownloadURLs(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)
	req := incomeRequest.(request.GetBatchMediaDownloadURLs)

	res, err := svc.GetBatchMediaDownloadURLs(ctx, req)
	return res, err
}

func VerifyMediaObject(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)
	req := incomeRequest.(request.VerifyMediaObject)

	res, err := svc.VerifyMediaObject(ctx, req)
	return res, err
}

func CreateFreeGroupChannel(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.CreateFreeGroupChannel)

	res, err := svc.CreateFreeGroupChannel(ctx, req)

	return res, err
}

func JoinFreeGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.JoinFreeGroup)

	err := svc.JoinFreeGroup(ctx, req)
	return nil, err
}

func GetFreeGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.GetFreeGroup)

	res, err := svc.GetFreeGroup(ctx, req.ChannelID)
	return res, err
}

func LeaveFreeGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.LeaveGroup)

	err := svc.LeaveFreeGroup(ctx, req)
	return nil, err
}

func UpdateFreeGroupMetadata(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.UpdateGroupMetadata)

	err := svc.UpdateFreeGroupMetadata(ctx, req)

	return nil, err
}

func UpdateFreeGroupLink(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.UpdateGroupLink)

	err := svc.UpdateFreeGroupLink(ctx, req)

	return nil, err
}

func ListUserFreeGroups(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.ListUserFreeGroups)

	res, err := svc.ListUserFreeGroups(ctx, req)
	return res, err
}

func DeleteFreeGroup(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.DeleteFreeGroup)

	err := svc.DeleteFreeGroup(ctx, req)
	return nil, err
}

func AddFreeGroupCoCreators(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.AddFreeGroupChatCoCreators)

	err := svc.AddFreeGroupCoCreators(ctx, req)
	return nil, err
}

func RemoveFreeGroupCoCreator(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(chat.Usecase)

	req := incomeRequest.(request.RemoveFreeGroupChatCoCreator)

	err := svc.RemoveFreeGroupCoCreator(ctx, req)
	return nil, err
}
