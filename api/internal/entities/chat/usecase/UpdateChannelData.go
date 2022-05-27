package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) UpdateChannelData(ctx context.Context, channel *chat.Channel, isStart bool, state string, isPublic bool) (*response.ChannelData, error) {
	requestURL := fmt.Sprintf(
		"https://api-%s.sendbird.com/v3/group_channels/%s",
		appconfig.Config.Sendbird.ApplicationID,
		channel.URL,
	)

	var expiresAt int64 = constants.GOAT_CHAT_EXPIRES_AT_DEFAULT_VALUE
	var newIsPublic bool

	if isStart {
		// Expire this channel 1 week from now in Unix milliseconds
		addYears := 0
		addMonths := 0
		addDays := 7
		expiresAt = time.Now().AddDate(addYears, addMonths, addDays).UnixMilli()

		newIsPublic = isPublic
	} else {
		oldIsPublic := true
		if channel.Data.IsConversationPublic != nil {
			oldIsPublic = *channel.Data.IsConversationPublic
		}
		newIsPublic = isPublic && oldIsPublic

		channel.Data.StartOfDraftMessageID = nil
		channel.Data.StartOfDraftMessageTS = nil
		channel.Data.EndOfDraftMessageTS = nil
	}

	// Update data field in SendBird
	data := channel.Data
	data.ExpiresAt = &expiresAt
	data.ChannelState = &state
	data.IsConversationPublic = &newIsPublic

	marshalledData, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		vlog.Errorf(ctx, "Error marshalling data struct: %v", marshalErr)
		return nil, marshalErr
	}

	putBody, jsonErr := json.Marshal(map[string]string{
		"data": string(marshalledData),
	})
	if jsonErr != nil {
		vlog.Errorf(ctx, "Error marshalling request map: %v", jsonErr)
		return nil, jsonErr
	}
	requestBody := bytes.NewBuffer(putBody)

	sendBirdReq, requestErr := http.NewRequest(
		"PUT",
		requestURL,
		requestBody,
	)

	if requestErr != nil {
		vlog.Errorf(ctx, "Something went wrong when creating request to update Sendbird channel data. Err: %v", requestErr)
		return nil, requestErr
	}

	sendBirdReq.Header.Set("Content-Type", "application/json")
	sendBirdReq.Header.Set("Api-Token", appconfig.Config.Sendbird.MasterAPIKey)

	sendBirdRes, sendBirdErr := svc.msg.Server.Do(sendBirdReq)
	if sendBirdErr != nil {
		return nil, sendBirdErr
	}
	if sendBirdRes == nil || sendBirdRes.StatusCode != http.StatusOK {
		vlog.Errorf(ctx, "Something went wrong when updating Sendbird channel data %s", channel.URL)
		return nil, fmt.Errorf("something went wrong updating Sendbird channel data")
	}

	return &data, nil
}
