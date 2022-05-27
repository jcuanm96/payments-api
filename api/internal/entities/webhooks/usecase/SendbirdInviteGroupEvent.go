package service

import (
	"context"
	"strconv"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) SendbirdInviteGroupEvent(ctx context.Context, event *sendbird.GroupChannelInviteEvent) error {
	if len(event.Invitees) == 0 || event.Channel.IsDistinct {
		return nil
	}

	for _, invitee := range event.Invitees {
		message, data, rangesErr := calculateInviteGroupAdminMessageRanges(ctx, &event.Inviter, &invitee)
		if rangesErr != nil {
			continue
		}
		svc.chat.SendGroupMessage(ctx, event.Channel.ChannelURL, *data.Ranges[0].ID, message, data, constants.ChannelEventCustomType)
	}
	return nil
}

func calculateInviteGroupAdminMessageRanges(ctx context.Context, inviter *sendbird.WebhookUser, invitee *sendbird.WebhookUser) (string, *response.AdminMessageData, error) {
	addedText := " added "
	message := inviter.NickName + addedText + invitee.NickName

	inviterID, atoiErr := strconv.Atoi(inviter.UserID)
	if atoiErr != nil {
		vlog.Errorf(ctx, "Error converting sendbird user ID %s to int: %s", inviter.UserID, atoiErr.Error())
		return "", nil, atoiErr
	}
	inviterNicknameRange := response.MessageRange{
		Start: 0,
		End:   utils.Utf16len(inviter.NickName),
		Type:  constants.AdminMessageRangeNickname,
		ID:    &inviterID,
	}

	textRange := response.MessageRange{
		Start: inviterNicknameRange.End,
		End:   inviterNicknameRange.End + utils.Utf16len(addedText),
		Type:  constants.AdminMessageRangeText,
	}

	inviteeID, atoiErr := strconv.Atoi(invitee.UserID)
	if atoiErr != nil {
		vlog.Errorf(ctx, "Error converting sendbird user ID %s to int: %s", invitee.UserID, atoiErr.Error())
		return "", nil, atoiErr
	}
	inviteeNicknameRange := response.MessageRange{
		Start: textRange.End,
		End:   utils.Utf16len(message),
		Type:  constants.AdminMessageRangeNickname,
		ID:    &inviteeID,
	}

	ranges := []response.MessageRange{inviterNicknameRange, textRange, inviteeNicknameRange}
	messageData := &response.AdminMessageData{
		Ranges: ranges,
	}

	return message, messageData, nil
}
