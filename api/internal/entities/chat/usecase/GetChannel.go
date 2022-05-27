package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

type SendBirdUser struct {
	UserID string `json:"user_id"`
}

type GetChannelResponse struct {
	Data        string                `json:"data"`
	Members     []SendBirdUser        `json:"members"`
	LastMessage *chat.SendBirdMessage `json:"last_message"`
}

func (svc *usecase) GetChannel(ctx context.Context, sendBirdChannelID string) (*chat.Channel, error) {
	requestURL := fmt.Sprintf(
		"https://api-%s.sendbird.com/v3/group_channels/%s?show_member=true",
		os.Getenv("SENDBIRD_APPLICATION_ID"),
		sendBirdChannelID,
	)

	sendBirdReq, requestErr := http.NewRequest(
		"GET",
		requestURL,
		nil,
	)

	if requestErr != nil {
		vlog.Errorf(ctx, "Something went wrong when creating request to check Sendbird channel membership. Err: %s", requestErr.Error())
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when checking channel membership",
		)
	}

	sendBirdReq.Header.Set("Content-Type", "application/json")
	sendBirdReq.Header.Set("Api-Token", os.Getenv("SENDBIRD_MASTER_API_KEY"))

	sendBirdRes, sendBirdErr := svc.msg.Server.Do(sendBirdReq)
	if sendBirdErr != nil {
		vlog.Errorf(ctx, "Something went wrong when getting Sendbird channel members for channel %s. Err: %s", sendBirdChannelID, sendBirdErr.Error())
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when getting channel members",
		)
	}
	if sendBirdRes == nil {
		vlog.Error(ctx, "Sendbird response for getting channel members was nil")
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when getting channel members",
		)
	}

	defer sendBirdRes.Body.Close()

	var getMembersRes GetChannelResponse
	json.NewDecoder(sendBirdRes.Body).Decode(&getMembersRes)

	// We're only supporting a creator chat between two users
	if len(getMembersRes.Members) != 2 {
		vlog.Errorf(ctx, "Error: did not get correct number of members when getting Sendbird channel members for channel id %s.", sendBirdChannelID)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when getting channel members",
		)
	}

	ids := map[int]struct{}{}
	for _, sendbirdUser := range getMembersRes.Members {
		userID, strConvErr := strconv.Atoi(sendbirdUser.UserID)
		if strConvErr != nil {
			vlog.Errorf(ctx, "Error converting string %s to int when getting Sendbird channel members for channel %s. Err: %s", sendbirdUser.UserID, sendBirdChannelID, strConvErr.Error())
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong when getting channel members",
			)
		}

		// Empty struct because we don't care about the value, only fast lookup for key
		ids[userID] = struct{}{}
	}

	channelData := response.ChannelData{}
	json.Unmarshal([]byte(getMembersRes.Data), &channelData)

	channel := chat.Channel{
		URL:           sendBirdChannelID,
		Data:          channelData,
		MemberUserIDs: ids,
		LastMessage:   getMembersRes.LastMessage,
	}

	return &channel, nil
}
