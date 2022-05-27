package webhooks

import (
	"context"
	"fmt"
	"strconv"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func StripeValidateSubscriptionMetadata(ctx context.Context, metadata map[string]string, subscriptionID string) (*wallet.StripeEventMetadata, error) {
	providerUserID, providerUserIDErr := wallet.GetIDFromMetadata(ctx, metadata, "providerUserID")
	if providerUserIDErr != nil {
		return nil, providerUserIDErr
	}
	customerUserID, customerUserIDErr := wallet.GetIDFromMetadata(ctx, metadata, "customerUserID")
	if customerUserIDErr != nil {
		return nil, customerUserIDErr
	}

	tierIDStr, tierOk := metadata["tierID"]
	var tierID *int
	if tierOk {
		tierIDValue, tierAtoiErr := strconv.Atoi(tierIDStr)
		if tierAtoiErr != nil {
			tierErr := fmt.Errorf("could not convert tierID %s to int in Stripe metadata for subscription %s. Err: %s", tierIDStr, subscriptionID, tierAtoiErr.Error())
			vlog.Errorf(ctx, "%v", tierErr)
			return nil, tierErr
		}
		tierID = &tierIDValue
	}

	var channelID *string
	channelIDValue, channelIDOk := metadata["sendbirdChannelID"]
	if channelIDOk {
		channelID = &channelIDValue
	}

	// isTrial cannot be set in the metadata because
	// then a subscription would permanently be a trial.
	subscriptionMetadata := wallet.StripeEventMetadata{
		ProviderUserID: providerUserID,
		CustomerUserID: customerUserID,
		TierID:         tierID,
		ChannelID:      channelID,
	}
	return &subscriptionMetadata, nil
}

// `messagef` should be of the form "%s some text here" where %s is the nickname of `user`
func CalculateNicknameThenTextRanges(ctx context.Context, user *sendbird.WebhookUser, messagef string) (string, *response.AdminMessageData, error) {
	message := fmt.Sprintf(messagef, user.NickName)

	userID, atoiErr := strconv.Atoi(user.UserID)
	if atoiErr != nil {
		vlog.Errorf(ctx, "Error converting sendbird user ID %s to int: %s", user.UserID, atoiErr.Error())
		return "", nil, atoiErr
	}
	nicknameRange := response.MessageRange{
		Start: 0,
		End:   utils.Utf16len(user.NickName),
		Type:  constants.AdminMessageRangeNickname,
		ID:    &userID,
	}

	textRange := response.MessageRange{
		Start: nicknameRange.End,
		End:   utils.Utf16len(message),
		Type:  constants.AdminMessageRangeText,
	}

	ranges := []response.MessageRange{nicknameRange, textRange}
	messageData := &response.AdminMessageData{
		Ranges: ranges,
	}

	return message, messageData, nil
}

func CalculateRemovedGroupAdminMessageRanges(banningUser *response.User, bannedUser *response.User) (string, *response.AdminMessageData, error) {
	banningUserNickname := fmt.Sprintf("%s %s", banningUser.FirstName, banningUser.LastName)
	bannedUserNickname := fmt.Sprintf("%s %s", bannedUser.FirstName, bannedUser.LastName)
	addedText := " removed "
	message := banningUserNickname + addedText + bannedUserNickname

	banningUserNicknameRange := response.MessageRange{
		Start: 0,
		End:   utils.Utf16len(banningUserNickname),
		Type:  constants.AdminMessageRangeNickname,
		ID:    &bannedUser.ID,
	}

	textRange := response.MessageRange{
		Start: banningUserNicknameRange.End,
		End:   banningUserNicknameRange.End + utils.Utf16len(addedText),
		Type:  constants.AdminMessageRangeText,
	}

	bannedUserNicknameRange := response.MessageRange{
		Start: textRange.End,
		End:   utils.Utf16len(message),
		Type:  constants.AdminMessageRangeNickname,
		ID:    &bannedUser.ID,
	}

	ranges := []response.MessageRange{banningUserNicknameRange, textRange, bannedUserNicknameRange}
	messageData := &response.AdminMessageData{
		Ranges: ranges,
	}

	return message, messageData, nil
}

func ConvertUserToWebhookUser(user *sendbird.User) sendbird.WebhookUser {
	return sendbird.WebhookUser{
		UserID:     user.UserID,
		NickName:   user.NickName,
		ProfileURL: user.ProfileURL,
	}
}
