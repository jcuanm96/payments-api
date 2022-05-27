package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

type IsMemberResponse struct {
	IsMember bool `json:"is_member"`
}

func (svc *usecase) IsUserInChannel(ctx context.Context, userID int, sendBirdChannelID string) (bool, error) {
	requestURL := fmt.Sprintf(
		"https://api-%s.sendbird.com/v3/group_channels/%s/members/%d",
		appconfig.Config.Sendbird.ApplicationID,
		sendBirdChannelID,
		userID)
	sendBirdReq, requestErr := http.NewRequest(
		"GET",
		requestURL,
		nil,
	)

	if requestErr != nil {
		vlog.Errorf(ctx, "Something went wrong when creating request to check Sendbird channel membership. Err: %s\n", requestErr.Error())
		return false, httperr.NewCtx(
			ctx, 500, http.StatusInternalServerError, "Error checking channel membership")
	}

	sendBirdReq.Header.Set("Content-Type", "application/json")
	sendBirdReq.Header.Set("Api-Token", appconfig.Config.Sendbird.MasterAPIKey)

	sendBirdRes, sendBirdErr := svc.msg.Server.Do(sendBirdReq)
	if sendBirdErr != nil || sendBirdRes == nil {
		vlog.Errorf(ctx, "Something went wrong when checking Sendbird channel membership Err: %s\n", sendBirdErr.Error())
		return false, httperr.NewCtx(
			ctx, 500, http.StatusInternalServerError, "Error checking channel membership")
	}

	defer sendBirdRes.Body.Close()

	var isMemberRes IsMemberResponse
	json.NewDecoder(sendBirdRes.Body).Decode(&isMemberRes)

	return isMemberRes.IsMember, nil
}
