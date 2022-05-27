package chat

import "github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"

type UpdateGoatChatResult struct {
	GoatChatMessagesID int
	IsPublic           bool
}

type SendBirdMessage struct {
	CreatedAt int64 `json:"created_at"`
}

type Channel struct {
	URL           string
	Data          response.ChannelData
	MemberUserIDs map[int]struct{}
	LastMessage   *SendBirdMessage
}
