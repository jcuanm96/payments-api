package mocks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

type MockSendBirdClient struct {
	httpClient *http.Client
	user       sendbird.User
}

func NewSendBirdMockClient() sendbird.Client {
	c := MockSendBirdClient{}
	c.httpClient = &http.Client{
		Timeout: time.Second * 5,
	}

	return &c
}

func (c *MockSendBirdClient) CreateUser(createUserReq *sendbird.CreateUserParams) (*sendbird.User, error) {
	user := sendbird.User{
		AccessToken: "fake-access-token-123",
	}

	c.user = user

	return &user, nil
}

func (c *MockSendBirdClient) GetGroupChannel(channelURL string, req sendbird.GetGroupChannelParams) (*sendbird.GroupChannel, error) {
	operator := sendbird.User{
		UserID: "1", // Creator should always be first user to be inserted into db
	}
	return &sendbird.GroupChannel{
		Name:        "Test Group Channel",
		ChannelURL:  channelURL,
		MemberCount: 100,
		Operators:   []sendbird.User{operator},
	}, nil
}

func (c *MockSendBirdClient) JoinGroupChannel(channelURL string, req *sendbird.JoinGroupChannelParams) error {
	return nil
}

func (c *MockSendBirdClient) LeaveGroupChannel(channelURL string, req *sendbird.LeaveGroupChannelParams) error {
	return nil
}

func (c *MockSendBirdClient) BlockUser(userID, blockUserID int) (*sendbird.User, error) {
	return nil, nil
}

func (c *MockSendBirdClient) UnblockUser(userID, unblockUserID int) error {
	return nil
}

func (c *MockSendBirdClient) GetUser(userID int) (*sendbird.User, error) {
	return &c.user, nil
}

func (c *MockSendBirdClient) UpsertUserMetadata(userID int, newMetadata interface{}) error {
	return nil
}

func (c *MockSendBirdClient) ListGroupChannelMessages(req sendbird.ListGroupChannelMessagesParams) (*sendbird.ListMessages, error) {
	return nil, nil
}

func (c *MockSendBirdClient) ListMyGroupChannels(userID int, req sendbird.ListMyGroupChannelsParams) ([]sendbird.GroupChannel, error) {
	return nil, nil
}

func (c *MockSendBirdClient) ListGroupChannels(req sendbird.ListGroupChannelsParams) ([]sendbird.GroupChannel, error) {
	return nil, nil
}

func (c *MockSendBirdClient) DeleteUser(userID int) error {
	return nil
}

func (c *MockSendBirdClient) SendMessage(channelURL string, body *sendbird.SendMessageParams) (*sendbird.Message, error) {
	return nil, nil
}

func (c *MockSendBirdClient) CreateGroupChannel(ctx context.Context, req *sendbird.CreateGroupChannelParams) (*sendbird.GroupChannel, error) {
	operator := sendbird.User{
		UserID: fmt.Sprint(req.OperatorIDs[0]),
	}
	return &sendbird.GroupChannel{
		Name:              req.Name,
		Operators:         []sendbird.User{operator},
		ChannelURL:        "sendbird_group_channel_123456789",
		MemberCount:       1,
		Members:           []sendbird.User{operator},
		JoinedMemberCount: 0,
	}, nil
}
func (c *MockSendBirdClient) UpdateGroupChannel(channelURL string, req *sendbird.UpdateGroupChannelParams, shouldUpsert bool) (*sendbird.GroupChannel, error) {
	return nil, nil
}
func (c *MockSendBirdClient) IsMemberInGroupChannel(channelURL string, userID string) (*sendbird.IsMemberInGroupChannel, error) {
	return nil, nil
}

func (c *MockSendBirdClient) ListGroupChannelMembers(channelURL string, req sendbird.ListGroupChannelMembersParams) (*sendbird.ListGroupChannelMembersResponse, error) {
	return nil, nil
}

func (c *MockSendBirdClient) DeleteGroupChannel(channelURL string) error {
	return nil
}

func (c *MockSendBirdClient) BanUserFromGroupChannel(channelURL string, req *sendbird.BanUserFromGroupChannelParams) error {
	return nil
}

func (c *MockSendBirdClient) UnbanUserFromGroupChannel(channelURL string, userID string) error {
	return nil
}

func (c *MockSendBirdClient) SearchMessages(req sendbird.SearchMessagesParams) (*sendbird.SearchMessagesResponse, error) {
	return nil, nil
}

func (c *MockSendBirdClient) ViewGroupChannelMessage(channelURL string, messageID string) (*sendbird.Message, error) {
	return nil, nil
}

func (c *MockSendBirdClient) UpdateUser(userID int, req *sendbird.UpdateUserParams) (*sendbird.User, error) {
	return nil, nil
}

func (c *MockSendBirdClient) InviteGroupChannelMembers(channelURL string, req *sendbird.InviteGroupChannelMembers) (*sendbird.GroupChannel, error) {
	return nil, nil
}
