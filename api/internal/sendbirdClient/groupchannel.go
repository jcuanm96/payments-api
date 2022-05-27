package sendbird

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/go-querystring/query"
)

type GroupChannel struct {
	Name                 string           `json:"name"`
	ChannelURL           string           `json:"channel_url"`
	CoverURL             string           `json:"cover_url"`
	CustomType           string           `json:"custom_type"`
	Data                 string           `json:"data"`
	IsDistinct           bool             `json:"is_distinct"`
	IsPublic             bool             `json:"is_public"`
	IsSuper              bool             `json:"is_super"`
	IsEphemeral          bool             `json:"is_ephemeral"`
	IsAccessCodeRequired bool             `json:"is_access_code_required"`
	MemberCount          int              `json:"member_count"`
	JoinedMemberCount    int              `json:"joined_member_count"`
	Members              []User           `json:"members"` // Only shows first 10 members
	Operators            []User           `json:"operators"`
	ReadReceipt          map[string]int64 `json:"read_receipt"`
	MaxMessageLength     int              `json:"max_length_message"`
	UnreadMessageCount   int              `json:"unread_message_count"`
	UnreadMentionCount   int              `json:"unread_mention_count"`
	LastMessage          LastMessage      `json:"last_message"`
	CreatedBy            User             `json:"created_by"`
	CreatedAt            int64            `json:"created_at"`
	Freeze               bool             `json:"freeze"`
}

type LastMessage struct {
	CreatedAt int64 `json:"created_at"`
	User      User  `json:"user"`
}

func (c *GroupChannel) HasOperator(userID string) bool {
	hasOperator := false
	for _, user := range c.Operators {
		if user.UserID == userID {
			hasOperator = true
			break
		}
	}
	return hasOperator
}

type GetGroupChannelParams struct {
	ShowDeliveryReceipt bool `json:"show_delivery_receipt" url:"show_delivery_receipt,omitempty"`
	ShowReadReceipt     bool `json:"show_read_receipt" url:"show_read_receipt,omitempty"`
	ShowMember          bool `json:"show_member" url:"show_member,omitempty"` // Only returns first 10 members
}

func (c *client) GetGroupChannel(channelURL string, req GetGroupChannelParams) (*GroupChannel, error) {
	pathString := fmt.Sprintf("/group_channels/%s", channelURL)
	parsedURL := c.PrepareUrl(pathString)

	queryValues, queryValuesErr := query.Values(req)
	if queryValuesErr != nil {
		return nil, queryValuesErr
	}

	rawQuery := queryValues.Encode()

	result := &GroupChannel{}
	getReqErr := c.get(parsedURL, rawQuery, result)

	if getReqErr != nil {
		switch getSendbirdUserErr := getReqErr.(type) {
		case SendbirdErrorResponse:
			if getSendbirdUserErr.Code == ResourceNotFound {
				return nil, ErrGroupChannelNotFound
			}
			return nil, getReqErr
		default:
			return nil, getReqErr
		}
	}

	return result, nil
}

type CreateGroupChannelParams struct {
	UserIDs     []int                 `json:"user_ids,omitempty"`
	Name        string                `json:"name"`
	CoverURL    string                `json:"cover_url,omitempty"`
	CoverFile   *multipart.FileHeader `json:"cover_file,omitempty"`
	CustomType  string                `json:"custom_type,omitempty"`
	Data        string                `json:"data"`
	IsDistinct  bool                  `json:"is_distinct"`
	IsPublic    bool                  `json:"is_public"`
	IsSuper     bool                  `json:"is_super"`
	OperatorIDs []int                 `json:"operator_ids,omitempty"`
	InviterID   string                `json:"inviter_id,omitempty"`
}

func (c *client) CreateGroupChannel(ctx context.Context, req *CreateGroupChannelParams) (*GroupChannel, error) {
	if req.CoverFile == nil {
		pathString := "/group_channels"
		parsedURL := c.PrepareUrl(pathString)

		result := GroupChannel{}
		postReqErr := c.post(parsedURL, req, &result)
		if postReqErr != nil {
			return nil, postReqErr
		}

		return &result, nil
	}
	return c.createGroupChannelWithFile(ctx, req)
}

type UpdateGroupChannelParams struct {
	Name        string   `json:"name,omitempty"`
	CustomType  string   `json:"custom_type,omitempty"`
	Data        string   `json:"data,omitempty"`
	IsDistinct  bool     `json:"is_distinct,omitempty"`
	OperatorIDs []string `json:"operator_ids,omitempty"`
}

func (c *client) UpdateGroupChannel(channelURL string, req *UpdateGroupChannelParams, shouldUpsertData bool) (*GroupChannel, error) {
	// Updating Sendbird Data is destructive, so we need to get the fields first and overwrite them with themselves
	if shouldUpsertData && req.Data != "" {
		channel, getChannelErr := c.GetGroupChannel(channelURL, GetGroupChannelParams{})
		if getChannelErr != nil {
			return nil, getChannelErr
		} else if channel == nil {
			return nil, errors.New("the specified Sendbird channel does not exist")
		}

		currData := map[string]interface{}{}
		unmarshalOldErr := json.Unmarshal([]byte(channel.Data), &currData)
		if unmarshalOldErr != nil {
			return nil, unmarshalOldErr
		}

		newData := map[string]interface{}{}
		unmarshalNewErr := json.Unmarshal([]byte(req.Data), &newData)
		if unmarshalNewErr != nil {
			return nil, unmarshalNewErr
		}

		// Merge the two data versions into one
		for key, value := range newData {
			currData[key] = value
		}

		newDataStr, marshalErr := json.Marshal(currData)
		if marshalErr != nil {
			return nil, marshalErr
		}

		req.Data = string(newDataStr)
	}

	pathString := fmt.Sprintf("/group_channels/%s", channelURL)
	parsedURL := c.PrepareUrl(pathString)

	result := GroupChannel{}
	putReqErr := c.put(parsedURL, req, &result)
	if putReqErr != nil {
		return nil, putReqErr
	}

	return &result, nil
}

func (c *client) IsMemberInGroupChannel(channelURL string, userID string) (*IsMemberInGroupChannel, error) {
	pathString := fmt.Sprintf("/group_channels/%s/members/%s", channelURL, userID)
	parsedURL := c.PrepareUrl(pathString)

	result := &IsMemberInGroupChannel{}
	getReqErr := c.get(parsedURL, "", result)
	if getReqErr != nil {
		return nil, getReqErr
	}

	return result, nil
}

type IsMemberInGroupChannel struct {
	IsMember bool `json:"is_member" `
}

type ListGroupChannelMembersParams struct {
	Token string `json:"token,omitempty" url:"token,omitempty"`
	Limit int    `json:"limit,omitempty" url:"limit,omitempty"`
}

type ListGroupChannelMembersResponse struct {
	Members []User `json:"members"`
	Next    string `json:"next"`
}

func (c *client) ListGroupChannelMembers(channelURL string, req ListGroupChannelMembersParams) (*ListGroupChannelMembersResponse, error) {
	pathString := fmt.Sprintf("/group_channels/%s/members", channelURL)
	parsedURL := c.PrepareUrl(pathString)

	queryValues, queryValuesErr := query.Values(req)
	if queryValuesErr != nil {
		return nil, queryValuesErr
	}

	rawQuery := queryValues.Encode()

	result := &ListGroupChannelMembersResponse{}
	getReqErr := c.get(parsedURL, rawQuery, result)
	if getReqErr != nil {
		return nil, getReqErr
	}

	return result, nil
}

func (c *client) DeleteGroupChannel(channelURL string) error {
	pathString := fmt.Sprintf("/group_channels/%s", channelURL)
	parsedURL := c.PrepareUrl(pathString)

	result := SendbirdErrorResponse{}
	deleteReqErr := c.delete(parsedURL, "", &result)
	if deleteReqErr != nil {
		return deleteReqErr
	}

	return nil
}

type ListGroupChannelsParams struct {
	MembersExactlyIn []string `json:"members_exactly_in,omitempty"` // manually encode this one
	ShowMember       *bool    `json:"show_member,omitempty" url:"show_member,omitempty"`
	DistinctMode     *string  `json:"distinct_mode,omitempty" url:"distinct_mode,omitempty"`
	PublicMode       *string  `json:"public_mode,omitempty" url:"public_mode,omitempty"`
	NameContains     *string  `json:"name_contains,omitempty" url:"name_contains,omitempty"`
	Limit            *int64   `json:"limit,omitempty" url:"limit,omitempty"`
	MembersIncludeIn []string `json:"members_include_in,omitempty"` // manually encode this one
}

type ListGroupChannelsResponse struct {
	Channels []GroupChannel `json:"channels"`
}

func (c *client) ListGroupChannels(req ListGroupChannelsParams) ([]GroupChannel, error) {
	pathString := "/group_channels"
	parsedURL := c.PrepareUrl(pathString)

	queryValues, queryValuesErr := query.Values(req)
	if queryValuesErr != nil {
		return nil, queryValuesErr
	}

	rawQuery := queryValues.Encode()
	rawQuery += "&members_exactly_in=" + strings.Join(req.MembersExactlyIn, ",")
	rawQuery += "&members_include_in=" + strings.Join(req.MembersIncludeIn, ",")

	response := ListGroupChannelsResponse{}

	getReqErr := c.get(parsedURL, rawQuery, &response)
	if getReqErr != nil {
		return nil, getReqErr
	}

	return response.Channels, nil
}

type ListMyGroupChannelsParams struct {
	CreatedBefore           int64  `json:"created_before" url:"created_before,omitempty"`
	Order                   string `json:"order" url:"order,omitempty"`
	Limit                   int64  `json:"limit" url:"limit,omitempty"`
	DistinctMode            string `json:"distinct_mode" url:"distinct_mode,omitempty"`
	NameContains            string `json:"name_contains,omitempty" url:"name_contains,omitempty"`
	MembersNicknameContains string `json:"members_nickname_contains,omitempty" url:"members_nickname_contains,omitempty"`
}

func (c *client) ListMyGroupChannels(userID int, req ListMyGroupChannelsParams) ([]GroupChannel, error) {
	pathString := fmt.Sprintf("/users/%d/my_group_channels", userID)
	parsedURL := c.PrepareUrl(pathString)

	queryValues, queryValuesErr := query.Values(req)
	if queryValuesErr != nil {
		return nil, queryValuesErr
	}

	rawQuery := queryValues.Encode()

	response := ListGroupChannelsResponse{}
	getReqErr := c.get(parsedURL, rawQuery, &response)
	if getReqErr != nil {
		return nil, getReqErr
	}

	return response.Channels, nil
}

type JoinGroupChannelParams struct {
	UserID string `json:"user_id"`
}

func (c *client) JoinGroupChannel(channelURL string, req *JoinGroupChannelParams) error {
	pathString := fmt.Sprintf("/group_channels/%s/join", channelURL)
	parsedURL := c.PrepareUrl(pathString)

	result := struct{}{}
	putReqErr := c.put(parsedURL, req, &result)
	if putReqErr != nil {
		return putReqErr
	}

	return nil
}

type LeaveGroupChannelParams struct {
	UserIDs []string `json:"user_ids"`
}

func (c *client) LeaveGroupChannel(channelURL string, req *LeaveGroupChannelParams) error {
	pathString := fmt.Sprintf("/group_channels/%s/leave", channelURL)
	parsedURL := c.PrepareUrl(pathString)

	result := struct{}{}
	putReqErr := c.put(parsedURL, req, &result)
	if putReqErr != nil {
		return putReqErr
	}

	return nil
}

type BanUserFromGroupChannelParams struct {
	UserID  string `json:"user_id"`
	AgentID string `json:"agent_id"`
}

func (c *client) BanUserFromGroupChannel(channelURL string, req *BanUserFromGroupChannelParams) error {
	pathString := fmt.Sprintf("/group_channels/%s/ban", channelURL)
	parsedURL := c.PrepareUrl(pathString)

	result := struct{}{}
	postReqErr := c.post(parsedURL, req, &result)
	if postReqErr != nil {
		return postReqErr
	}

	return nil
}

func (c *client) UnbanUserFromGroupChannel(channelURL string, userID string) error {
	pathString := fmt.Sprintf("/group_channels/%s/ban/%s", channelURL, userID)
	parsedURL := c.PrepareUrl(pathString)

	result := SendbirdErrorResponse{}
	deleteReqErr := c.delete(parsedURL, "", &result)
	if deleteReqErr != nil {
		return deleteReqErr
	}

	return nil
}

type InviteGroupChannelMembers struct {
	UserIDs []string `json:"user_ids,omitempty"`
	Users   []string `json:"users,omitempty"`
}

func (c *client) InviteGroupChannelMembers(channelURL string, req *InviteGroupChannelMembers) (*GroupChannel, error) {
	pathString := fmt.Sprintf("/group_channels/%s/invite", channelURL)
	parsedURL := c.PrepareUrl(pathString)

	result := &GroupChannel{}
	postReqErr := c.post(parsedURL, req, result)
	if postReqErr != nil {
		return nil, postReqErr
	}

	return result, nil
}
