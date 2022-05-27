package response

import (
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

type SearchGlobal struct {
	GlobalChannels []sendbird.GroupChannel `json:"globalChannels"`
	MyChannels     []sendbird.GroupChannel `json:"myChannels"`
	Users          []User                  `json:"users"`
	Messages       []SearchMessage         `json:"messages"`
}

type SearchMessage struct {
	Message    sendbird.Message `json:"message"`
	StartRange int              `json:"startRange"`
	EndRange   int              `json:"endRange"`
}

type SearchMention struct {
	Users []sendbird.User `json:"user"`
}
