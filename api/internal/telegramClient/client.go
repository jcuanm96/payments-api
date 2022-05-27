package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/sirupsen/logrus"
)

var TelegramClient Client

type Client interface {
	SendMessage(text string)
}

type client struct {
	httpClient *http.Client
	apiToken   string
}

type TelegramErrorResponse struct {
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func Init() {
	c := client{}
	c.httpClient = &http.Client{
		Timeout: time.Second * 5,
	}
	c.apiToken = appconfig.Config.Telegram.ApiToken

	TelegramClient = &c
}

func (c *client) PrepareUrl(chatID string, text string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", c.apiToken, chatID, text)
}

func (c *client) SendMessage(text string) {
	if appconfig.Config.Gcloud.Project != "vama-prod" && appconfig.Config.Gcloud.Project != "vama-staging" {
		return
	}

	text = fmt.Sprintf("%s: %s", appconfig.Config.Gcloud.Project, text)
	url := c.PrepareUrl(appconfig.Config.Telegram.ChannelIDs.VamaAlerts, url.QueryEscape(text))
	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		logrus.Errorf("Error creating request in telegram client. Err: %s", reqErr.Error())
		return
	}

	httpResp, getHttpErr := c.httpClient.Do(req)
	if getHttpErr != nil {
		logrus.Errorf("Error executing GET request in telegram client. Err: %s", getHttpErr.Error())
		return
	}

	defer httpResp.Body.Close()

	telegramErrRes := TelegramErrorResponse{}
	decodeErr := json.NewDecoder(httpResp.Body).Decode(&telegramErrRes)
	if decodeErr != nil {
		logrus.Errorf("Error decoding response body in telegram client. Err: %s", decodeErr.Error())
		return
	}

	if httpResp.StatusCode != 200 {
		logrus.Errorf("Error the telegram client returned a non-success code %d. Message: %v", httpResp.StatusCode, telegramErrRes.Description)
		return
	}
}
