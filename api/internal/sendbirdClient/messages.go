package sendbird

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

const AdminMessageType = "ADMM"
const FileMessageType = "FILE"
const TextMessageType = "MESG"

type Message struct {
	MessageSurvivalSeconds int64         `json:"message_survival_seconds"`
	CustomType             string        `json:"custom_type"`
	MentionedUsers         []User        `json:"mentioned_users"`
	Translations           *Translations `json:"translations,omitempty"`
	UpdatedAt              int64         `json:"updated_at"`
	IsOpMsg                bool          `json:"is_op_msg"`
	IsRemoved              bool          `json:"is_removed"`
	User                   User          `json:"user"`
	File                   File          `json:"file"`
	Message                *string       `json:"message,omitempty"`
	Data                   *string       `json:"data,omitempty"`
	MessageRetentionHour   int64         `json:"message_retention_hour"`
	Silent                 bool          `json:"silent"`
	Type                   string        `json:"type"`
	CreatedAt              int64         `json:"created_at"`
	ChannelType            string        `json:"channel_type"`
	ReqID                  string        `json:"req_id"`
	MentionType            string        `json:"mention_type"`
	ChannelURL             string        `json:"channel_url"`
	MessageID              int64         `json:"message_id"`
	Thumbnails             []interface{} `json:"thumbnails,omitempty"`
	RequireAuth            *bool         `json:"require_auth,omitempty"`
}

type File struct {
	URL  *string `json:"url,omitempty"`
	Data *string `json:"data,omitempty"`
	Type *string `json:"type,omitempty"`
	Name *string `json:"name,omitempty"`
	Size *int64  `json:"size,omitempty"`
}

type Translations struct {
}

type ListGroupChannelMessagesParams struct {
	ChannelURL string  `json:"channel_url" url:"channel_url"`
	MessageTS  int64   `json:"message_ts" url:"message_ts"`
	PrevLimit  *int    `json:"prev_limit,omitempty" url:"prev_limit,omitempty"`
	NextLimit  *int    `json:"next_limit,omitempty" url:"next_limit,omitempty"`
	Reverse    *bool   `json:"reverse,omitempty" url:"reverse,omitempty"`
	SenderID   *string `json:"sender_id,omitempty" url:"sender_id,omitempty"`
	Include    *bool   `json:"include,omitempty" url:"include,omitempty"`
}

type ListMessages struct {
	Messages []Message `json:"messages"`
}

func (c *client) ListGroupChannelMessages(req ListGroupChannelMessagesParams) (*ListMessages, error) {
	pathString := fmt.Sprintf("/group_channels/%s/messages", req.ChannelURL)
	parsedURL := c.PrepareUrl(pathString)

	queryValues, queryValuesErr := query.Values(req)
	if queryValuesErr != nil {
		return nil, queryValuesErr
	}
	rawQuery := queryValues.Encode()

	result := ListMessages{}

	getReqErr := c.get(parsedURL, rawQuery, &result)
	if getReqErr != nil {
		return nil, getReqErr
	}

	return &result, nil
}

type SendMessageParams struct {
	MessageType string  `json:"message_type"`
	UserID      string  `json:"user_id"`
	Message     string  `json:"message"`
	Data        *string `json:"data,omitempty"`
	CustomType  string  `json:"custom_type"`
	IsSilent    bool    `json:"is_silent"`
	SendPush    *bool   `json:"send_push,omitempty"`
}

func (c *client) SendMessage(channelURL string, body *SendMessageParams) (*Message, error) {
	pathString := fmt.Sprintf("/group_channels/%s/messages", channelURL)
	parsedURL := c.PrepareUrl(pathString)

	result := Message{}
	postReqErr := c.post(parsedURL, body, &result)
	if postReqErr != nil {
		return nil, postReqErr
	}

	return &result, nil
}

type SearchMessagesParams struct {
	UserID string `json:"user_id" url:"user_id"`
	Query  string `json:"query" url:"query"`
	Limit  int    `json:"limit" url:"limit"`
}

type SearchMessagesResponse struct {
	TotalCount int       `json:"total_count"`
	Results    []Message `json:"results"`
}

func (c *client) SearchMessages(req SearchMessagesParams) (*SearchMessagesResponse, error) {
	pathString := "/search/messages"
	parsedURL := c.PrepareUrl(pathString)

	queryValues, queryValuesErr := query.Values(req)
	if queryValuesErr != nil {
		return nil, queryValuesErr
	}

	rawQuery := queryValues.Encode()

	response := &SearchMessagesResponse{}
	getReqErr := c.get(parsedURL, rawQuery, response)
	if getReqErr != nil {
		return nil, getReqErr
	}

	return response, nil
}

func (c *client) ViewGroupChannelMessage(channelURL string, messageID string) (*Message, error) {
	pathString := fmt.Sprintf("/group_channels/%s/messages/%s", channelURL, messageID)
	parsedURL := c.PrepareUrl(pathString)

	result := Message{}

	getReqErr := c.get(parsedURL, "", &result)
	if getReqErr != nil {
		return nil, getReqErr
	}

	return &result, nil
}
