package search

import sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"

type Channels struct {
	GlobalChannels []sendbird.GroupChannel
	MyChannels     []sendbird.GroupChannel
}
