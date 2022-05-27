package test

import (
	"fmt"
	"testing"

	"github.com/VamaSingapore/vama-api/cmd/test"
	test_util "github.com/VamaSingapore/vama-api/cmd/test/util"
	test_data "github.com/VamaSingapore/vama-api/cmd/test/util/data"
	"github.com/VamaSingapore/vama-api/internal/controller"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/gofiber/fiber/v2"
	"github.com/houqp/gtest"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ChatTests struct {
	app *fiber.App
	ctr *controller.Ctr
	db  *pgxpool.Pool
}

func (m *ChatTests) Setup(t *testing.T) {
	fmt.Println("Starting test server for Chat tests")
	app := test.StartTestServer()
	m.app = app.App
	m.db = app.Db
	m.ctr = app.Ctr
	test_util.InitializeDbSchemas(m.db)
}

func (m *ChatTests) Teardown(t *testing.T) {
	fmt.Println("Killing test server for Chat tests")
	m.app.Shutdown()
	m.db.Close()
}

func (m *ChatTests) BeforeEach(t *testing.T) {
	test_util.ClearDB(m.db)
	test_data.FillBootstrapUserData(m.db)
}

func (m *ChatTests) AfterEach(t *testing.T) {
}

// Test create paid group chat returns proper response
func (m *ChatTests) SubTestCreatePaidGroupChannel(t *testing.T) {
	testGroupName := "Test Paid Group Chat 1"
	params := map[string]interface{}{
		"priceInSmallestDenom": 300,
		"currency":             "usd",
		"name":                 testGroupName,
		"linkSuffix":           "testpaidgroup1",
	}

	asUserID := 1
	endpoint := "/api/v1/chat/paid/group/create"
	resp := &response.CreatePaidGroupChannel{}
	test_util.MakePostRequestAssert200(t, m.app, endpoint, params, &asUserID, resp)

	assert.Equal(t, testGroupName, resp.Channel.Name)
	require.Equal(t, 1, len(resp.Channel.Operators))
	assert.Equal(t, fmt.Sprint(asUserID), resp.Channel.Operators[0].UserID)
	assert.Equal(t, 1, resp.Channel.MemberCount)
	assert.Equal(t, "sendbird_group_channel_123456789", resp.Channel.ChannelURL)
	assert.Equal(t, 0, resp.Channel.JoinedMemberCount)
	require.Equal(t, 1, len(resp.Channel.Members))
	assert.Equal(t, fmt.Sprint(asUserID), resp.Channel.Members[0].UserID)
}

// Tests join paid group chat returns proper response
func (m *ChatTests) SubTestJoinPaidGroup(t *testing.T) {
	test_data.FillBootstrapCustomerUserSignupData(t, m.app)
	sendbirdChannelID := "sendbird_group_channel_123456789"
	params := map[string]interface{}{
		"channelID": sendbirdChannelID,
	}

	asUserID := 2
	asCreatorID := 1
	test_data.FillBootstrapPaidGroupChatData(m.db, m.ctr, asCreatorID, sendbirdChannelID, false)
	endpoint := "/api/v1/chat/paid/group/join"
	resp := &response.PaidGroupChatSubscription{}
	test_util.MakePostRequestAssert200(t, m.app, endpoint, params, &asUserID, resp)

	// ID and Group are only set when we're getting a paid group, not when we're upserting
	assert.Equal(t, 0, resp.ID)
	assert.Nil(t, resp.Group)

	assert.Equal(t, asUserID, resp.UserID)
	require.NotNil(t, resp.GoatUser)
	assert.Equal(t, asCreatorID, resp.GoatUser.ID)
	assert.True(t, resp.IsRenewing)
}

// Test member limit denies entry into paid group chat after reaching capacity
// and that a 403 error code was thrown
func (m *ChatTests) SubTestJoinGroupChatFailsFromMemberLimitThreshold(t *testing.T) {
	// Create paid group with member limit
	testGroupName := "Test Paid Group Chat 1"
	createPaidGroupParams := map[string]interface{}{
		"priceInSmallestDenom": 300,
		"currency":             "usd",
		"name":                 testGroupName,
		"linkSuffix":           "testpaidgroup1",
		"isMemberLimitEnabled": true,
	}

	asUserCreatorID := 1
	createEndpoint := "/api/v1/chat/paid/group/create"
	createPaidGroupResp := &response.CreatePaidGroupChannel{}
	test_util.MakePostRequestAssert200(t, m.app, createEndpoint, createPaidGroupParams, &asUserCreatorID, createPaidGroupResp)

	// Have user try to join at capacity paid group
	test_data.FillBootstrapCustomerUserSignupData(t, m.app)
	sendbirdChannelID := "sendbird_group_channel_123456789"
	joinPaidGroupParams := map[string]interface{}{
		"channelID": sendbirdChannelID,
	}

	asCustomerID := 2
	joinEndpoint := "/api/v1/chat/paid/group/join"
	joinPaidGroupResp := &response.PaidGroupChatSubscription{}
	test_util.MakePostRequestAssert403(t, m.app, joinEndpoint, joinPaidGroupParams, &asCustomerID, joinPaidGroupResp)
}

// Tests leave paid group chat returns proper response
func (m *ChatTests) SubTestLeavePaidGroupChat(t *testing.T) {
	test_data.FillBootstrapPaidGroupChatWithMembers(t, m.app, m.db, m.ctr)

	sendbirdChannelID := "sendbird_group_channel_123456789"
	leavePaidGroupParams := map[string]interface{}{
		"channelID": sendbirdChannelID,
	}

	asCustomerID := 2
	leaveEndpoint := "/api/v1/chat/paid/group/leave"
	leavePaidGroupResp := &response.LeavePaidGroup{}
	test_util.MakePostRequestAssert200(t, m.app, leaveEndpoint, leavePaidGroupParams, &asCustomerID, leavePaidGroupResp)

	// Group is only set when we're getting a paid group, not when we're upserting
	assert.Nil(t, leavePaidGroupResp.Subscription.Group)

	require.NotNil(t, leavePaidGroupResp.Subscription)
	assert.Equal(t, 1, leavePaidGroupResp.Subscription.ID)
	assert.Equal(t, asCustomerID, leavePaidGroupResp.Subscription.UserID)
	assert.Equal(t, "", leavePaidGroupResp.Subscription.StripeSubscriptionID)
	require.NotNil(t, leavePaidGroupResp.Subscription.GoatUser)
	assert.Equal(t, 1, leavePaidGroupResp.Subscription.GoatUser.ID)
	assert.False(t, leavePaidGroupResp.Subscription.IsRenewing)
}

// Tests cancel paid group chat returns proper response
func (m *ChatTests) SubTestCancelPaidGroupChat(t *testing.T) {
	test_data.FillBootstrapPaidGroupChatWithMembers(t, m.app, m.db, m.ctr)

	sendbirdChannelID := "sendbird_group_channel_123456789"
	leavePaidGroupParams := map[string]interface{}{
		"channelID": sendbirdChannelID,
	}

	asCustomerID := 2
	cancelEndpoint := "/api/v1/chat/paid/group/cancel"
	cancelPaidGroupResp := &response.CancelPaidGroup{}
	test_util.MakePostRequestAssert200(t, m.app, cancelEndpoint, leavePaidGroupParams, &asCustomerID, cancelPaidGroupResp)

	// Group is only set when we're getting a paid group, not when we're upserting
	assert.Nil(t, cancelPaidGroupResp.Subscription.Group)

	require.NotNil(t, cancelPaidGroupResp.Subscription)
	assert.Equal(t, 1, cancelPaidGroupResp.Subscription.ID)
	assert.Equal(t, asCustomerID, cancelPaidGroupResp.Subscription.UserID)
	assert.Equal(t, "", cancelPaidGroupResp.Subscription.StripeSubscriptionID)
	require.NotNil(t, cancelPaidGroupResp.Subscription.GoatUser)
	assert.Equal(t, 1, cancelPaidGroupResp.Subscription.GoatUser.ID)
	assert.False(t, cancelPaidGroupResp.Subscription.IsRenewing)
}

// Tests delete paid group chat returns proper response
func (m *ChatTests) SubTestDeletePaidGroupChat(t *testing.T) {
	test_data.FillBootstrapPaidGroupChatWithMembers(t, m.app, m.db, m.ctr)

	sendbirdChannelID := "sendbird_group_channel_123456789"
	deletePaidGroupParams := map[string]interface{}{
		"channelID": sendbirdChannelID,
	}

	asCreatorID := 1
	deleteEndpoint := "/api/v1/chat/paid/group"
	deleteResp := struct{}{}
	test_util.MakeDeleteRequestAssert200(t, m.app, deleteEndpoint, deletePaidGroupParams, asCreatorID, &deleteResp)
}

func TestChat(t *testing.T) {
	gtest.RunSubTests(t, &ChatTests{})
}
