package vamabot

import (
	"context"

	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

type Client interface {
	SendWelcomeMessages(ctx context.Context, userID int) error
}

type client struct {
	vamaUserID     int
	sendbirdClient sendbird.Client
}

func NewClient(vamaUserID int, sendbirdClient sendbird.Client) Client {
	return &client{
		vamaUserID:     vamaUserID,
		sendbirdClient: sendbirdClient,
	}
}
