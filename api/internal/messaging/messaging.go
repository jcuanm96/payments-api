package messaging

import (
	"net/http"
	"time"
)

type Messenger interface{}

type Client struct {
	Server *http.Client
}

func New() *Client {
	sendBirdClient := &http.Client{
		Timeout: time.Second * 10,
	}

	return &Client{sendBirdClient}
}
