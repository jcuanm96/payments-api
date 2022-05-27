package test

import (
	"fmt"
	"testing"

	"github.com/VamaSingapore/vama-api/cmd/test"
	test_util "github.com/VamaSingapore/vama-api/cmd/test/util"
	test_data "github.com/VamaSingapore/vama-api/cmd/test/util/data"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/gofiber/fiber/v2"
	"github.com/houqp/gtest"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FollowTests struct {
	app *fiber.App
	db  *pgxpool.Pool
}

func (m *FollowTests) Setup(t *testing.T) {
	fmt.Println("Starting test server for Follow tests")
	app := test.StartTestServer()
	m.app = app.App
	m.db = app.Db
	test_util.InitializeDbSchemas(m.db)
}

func (m *FollowTests) Teardown(t *testing.T) {
	fmt.Println("Killing test server for Follow tests")
	m.app.Shutdown()
	m.db.Close()
}

func (m *FollowTests) BeforeEach(t *testing.T) {
	test_util.ClearDB(m.db)
	test_data.FillBootstrapUserData(m.db)
}

func (m *FollowTests) AfterEach(t *testing.T) {
}

// Test user2 following user1 returns a 200.
func (m *FollowTests) SubTestFollow(t *testing.T) {
	params := map[string]interface{}{
		"userID": 1,
	}
	asUserID := 2
	resp := &response.User{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/follows/follow", params, &asUserID, resp)

	assert.Equal(t, 1, resp.ID)
}

// Test getting user2's following returns a 200
func (m *FollowTests) SubTestGetFollows(t *testing.T) {
	test_data.FillBootstrapFollowData(m.db)
	params := map[string]string{}
	asUserID := 2
	resp := &response.GetFollowedGoats{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/follows/goats", params, asUserID, resp)

	require.NotNil(t, resp)
	require.Equal(t, 1, len(resp.Goats))
	assert.Equal(t, 1, resp.Goats[0].ID)
}

// Test checking if user2 follows user1 returns a 200
func (m *FollowTests) SubTestIsFollow(t *testing.T) {
	test_data.FillBootstrapFollowData(m.db)
	params := map[string]string{
		"userID": "1",
	}
	asUserID := 2
	resp := &response.IsFollowing{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/follows/check", params, asUserID, resp)

	require.NotNil(t, resp)
	assert.Equal(t, true, resp.IsFollowing)
}

// Test user2 unfollowing user1 returns a 200
func (m *FollowTests) SubTestUnFollow(t *testing.T) {
	test_data.FillBootstrapFollowData(m.db)
	params := map[string]interface{}{
		"userID": 1,
	}
	asUserID := 1
	resp := &struct{}{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/follows/unfollow", params, &asUserID, resp)
}

// Test user3 follows user1, then unfollows, then checks if following.
func (m *FollowTests) SubTestFollowUnfollow(t *testing.T) {
	asUserID := 3
	followParams := map[string]interface{}{
		"userID": 1,
	}
	followResp := &response.User{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/follows/follow", followParams, &asUserID, followResp)
	assert.Equal(t, 1, followResp.ID)

	unfollowParams := map[string]interface{}{
		"userID": 1,
	}
	unfollowResp := &struct{}{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/follows/unfollow", unfollowParams, &asUserID, unfollowResp)

	isFollowingParams := map[string]string{
		"userID": "1",
	}
	isFollowingResp := &response.IsFollowing{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/follows/check", isFollowingParams, asUserID, isFollowingResp)
	require.NotNil(t, isFollowingResp)
	assert.Equal(t, false, isFollowingResp.IsFollowing)
}

func TestFollow(t *testing.T) {
	gtest.RunSubTests(t, &FollowTests{})
}
