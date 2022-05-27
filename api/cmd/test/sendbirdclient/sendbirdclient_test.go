package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/VamaSingapore/vama-api/cmd/test"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/gofiber/fiber/v2"
	"github.com/houqp/gtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sendbirdClientTests struct {
	client sendbird.Client
	app    *fiber.App
}

func (m *sendbirdClientTests) Setup(t *testing.T) {
	fmt.Println("Starting test server for SendbirdClient tests")
	// Need to start server to initialize appconfig for sendbird creds
	app := test.StartTestServer()
	m.app = app.App
	m.client = sendbird.NewClient()
}

func (m *sendbirdClientTests) Teardown(t *testing.T) {
	fmt.Println("Killing test server for SendbirdClient tests")
	m.app.Shutdown()
}
func (m *sendbirdClientTests) BeforeEach(t *testing.T) {}

func (m *sendbirdClientTests) AfterEach(t *testing.T) {}

// func (m *sendbirdClientTests) SubTestCreateThenDeleteUser(t *testing.T) {
// 	userID := 1
// 	user, createUserErr := m.client.CreateUser(&sendbird.CreateUserParams{
// 		UserID:           userID,
// 		Nickname:         "user1first user1last",
// 		ProfileURL:       "",
// 		IssueAccessToken: true,
// 	})

// 	require.Nil(t, createUserErr)
// 	assert.Equal(t, fmt.Sprint(userID), user.UserID)

// 	deleteUserErr := m.client.DeleteUser(userID)
// 	require.Nil(t, deleteUserErr)
// }

// func (m *sendbirdClientTests) SubTestCreateThenDeleteChannel(t *testing.T) {
// 	userID1 := 1
// 	user1, createUserErr := m.client.CreateUser(&sendbird.CreateUserParams{
// 		UserID:           userID1,
// 		Nickname:         "user1first user1last",
// 		ProfileURL:       "",
// 		IssueAccessToken: true,
// 	})
// 	require.Nil(t, createUserErr)
// 	assert.Equal(t, fmt.Sprint(userID1), user1.UserID)

// 	userID2 := 2
// 	user2, createUserErr := m.client.CreateUser(&sendbird.CreateUserParams{
// 		UserID:           userID2,
// 		Nickname:         "user2first user2last",
// 		ProfileURL:       "",
// 		IssueAccessToken: true,
// 	})
// 	require.Nil(t, createUserErr)
// 	assert.Equal(t, fmt.Sprint(userID2), user2.UserID)

// 	groupChannel, createChannelErr := m.client.CreateGroupChannel(&sendbird.CreateGroupChannelParams{
// 		UserIDs: []int{userID1, userID2},
// 		Data:    "{}",
// 	})
// 	require.Nil(t, createChannelErr)

// 	deleteChannelErr := m.client.DeleteGroupChannel(groupChannel.ChannelURL)
// 	require.Nil(t, deleteChannelErr)

// 	deleteUser1Err := m.client.DeleteUser(userID1)
// 	require.Nil(t, deleteUser1Err)

// 	deleteUser2Err := m.client.DeleteUser(userID2)
// 	require.Nil(t, deleteUser2Err)
// }

func (m *sendbirdClientTests) SubTestSendMessageThenListMessage(t *testing.T) {
	ctx := context.Background()
	userID1 := 1
	userID1String := fmt.Sprint(userID1)
	user1, createUserErr := m.client.CreateUser(&sendbird.CreateUserParams{
		UserID:           userID1,
		Nickname:         "user1first user1last",
		ProfileURL:       "",
		IssueAccessToken: true,
	})
	require.Nil(t, createUserErr)
	assert.Equal(t, userID1String, user1.UserID)

	userID2 := 2
	userID2String := fmt.Sprint(userID2)
	user2, createUserErr := m.client.CreateUser(&sendbird.CreateUserParams{
		UserID:           userID2,
		Nickname:         "user2first user2last",
		ProfileURL:       "",
		IssueAccessToken: true,
	})
	require.Nil(t, createUserErr)
	assert.Equal(t, userID2String, user2.UserID)

	groupChannel, createChannelErr := m.client.CreateGroupChannel(ctx, &sendbird.CreateGroupChannelParams{
		UserIDs: []int{userID1, userID2},
		Data:    "{}",
	})
	require.Nil(t, createChannelErr)

	mesgContent1 := "Hey this is the first test message!"
	mesgContent2 := "This is message number 2!"

	sendBody := &sendbird.SendMessageParams{
		MessageType: "MESG",
		UserID:      userID1String,
		Message:     mesgContent1,
	}

	message1, sendMessageErr := m.client.SendMessage(groupChannel.ChannelURL, sendBody)
	require.Nil(t, sendMessageErr)
	require.NotNil(t, message1.Message)
	assert.Equal(t, mesgContent1, *message1.Message)

	sendBody = &sendbird.SendMessageParams{
		MessageType: "MESG",
		UserID:      userID2String,
		Message:     mesgContent2,
	}

	message2, sendMessageErr := m.client.SendMessage(groupChannel.ChannelURL, sendBody)
	require.Nil(t, sendMessageErr)
	require.NotNil(t, message2.Message)
	assert.Equal(t, mesgContent2, *message2.Message)

	listMessagesParams := sendbird.ListGroupChannelMessagesParams{
		ChannelURL: groupChannel.ChannelURL,
		MessageTS:  message1.CreatedAt,
	}
	listMessages, listMessagesErr := m.client.ListGroupChannelMessages(listMessagesParams)
	require.Nil(t, listMessagesErr)
	messages := listMessages.Messages
	assert.Equal(t, 2, len(messages))
	require.NotNil(t, messages[0].Message)
	require.NotNil(t, messages[1].Message)
	assert.Equal(t, mesgContent1, *messages[0].Message)
	assert.Equal(t, mesgContent2, *messages[1].Message)

	deleteChannelErr := m.client.DeleteGroupChannel(groupChannel.ChannelURL)
	require.Nil(t, deleteChannelErr)

	deleteUser1Err := m.client.DeleteUser(userID1)
	require.Nil(t, deleteUser1Err)

	deleteUser2Err := m.client.DeleteUser(userID2)
	require.Nil(t, deleteUser2Err)
}

func TestSendbirdClient(t *testing.T) {
	gtest.RunSubTests(t, &sendbirdClientTests{})
}
