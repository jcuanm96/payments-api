package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/VamaSingapore/vama-api/cmd/test"
	test_util "github.com/VamaSingapore/vama-api/cmd/test/util"
	test_data "github.com/VamaSingapore/vama-api/cmd/test/util/data"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	userrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/gofiber/fiber/v2"
	"github.com/houqp/gtest"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FeedTests struct {
	app *fiber.App
	db  *pgxpool.Pool
}

func (m *FeedTests) Setup(t *testing.T) {
	fmt.Println("Starting test server for Feed tests")
	app := test.StartTestServer()
	m.app = app.App
	m.db = app.Db
	test_util.InitializeDbSchemas(m.db)
}

func (m *FeedTests) Teardown(t *testing.T) {
	fmt.Println("Killing test server for Feed tests")
	m.app.Shutdown()
	m.db.Close()
}

func (m *FeedTests) BeforeEach(t *testing.T) {
	test_util.ClearDB(m.db)
	test_data.FillBootstrapUserData(m.db)
	test_data.FillBootstrapFeedPostData(m.db)
}

func (m *FeedTests) AfterEach(t *testing.T) {
}

// Test getting feed post comments /feed/posts/comments
// Checks bootstrap post and comments are gettable
// and return a 200.
func (m *FeedTests) SubTestGetFeedPostComments(t *testing.T) {
	params := map[string]string{
		"postID":            "1",
		"limit":             "10",
		"lastUserCommentID": "0",
	}
	asUserID := 1
	resp := &response.GetFeedPostComments{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/feed/posts/comments", params, asUserID, resp)

	assert.Equal(t, 2, len(resp.Comments))
	assert.Equal(t, 1, resp.Comments[0].PostID)
	// Assert that comments are newest to oldest
	assert.Equal(t, "This is my second test comment on my first test feed post!", resp.Comments[0].TextContent)
	assert.Equal(t, "This is my first test comment on my first test feed post!", resp.Comments[1].TextContent)
}

func (m *FeedTests) SubTestMakeComment(t *testing.T) {
	params := map[string]interface{}{
		"postID": 1,
		"text":   "Yo Yo Yo testy tester in the house doing some testin",
	}
	asUserID := 1
	resp := &response.MakeComment{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/feed/posts/comments", params, &asUserID, resp)
}

// Tests that a blocked user cannot comment on a feed post.
func (m *FeedTests) SubTestMakeCommentFromBlockedUser(t *testing.T) {
	// creator user 1 blocks user 2
	ctx := context.Background()
	blockErr := userrepo.BlockUser(ctx, m.db, 1, 2)
	if blockErr != nil {
		vlog.Fatalf(ctx, "Error blocking user: %s", blockErr)
	}

	params := map[string]interface{}{
		"postID": 1,
		"text":   "Trying to comment but i'm blocked hehe",
	}
	asUserID := 2
	resp := &response.MakeComment{}
	expectedStatusCode := 400
	test_util.MakePostRequest(t, m.app, "/api/v1/feed/posts/comments", expectedStatusCode, params, &asUserID, resp)
}

func (m *FeedTests) SubTestDeleteComment(t *testing.T) {
	params := map[string]interface{}{
		"commentID": 1,
	}
	asUserID := 1
	resp := &map[string]string{}
	test_util.MakeDeleteRequestAssert200(t, m.app, "/api/v1/feed/comments", params, asUserID, resp)
}

func (m *FeedTests) SubTestGetFeedPostByID(t *testing.T) {
	params := map[string]string{
		"postID": "1",
	}

	asUserID := 1
	resp := &response.FeedPost{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/feed/posts", params, asUserID, resp)

	require.NotNil(t, resp.Creator)
	assert.Equal(t, 1, resp.Creator.ID)
	assert.Nil(t, resp.Customer)
}

func (m *FeedTests) SubTestGetFeedPostByLinkSuffix(t *testing.T) {
	params := map[string]string{
		"linkSuffix": "Gg4123p893q",
	}

	asUserID := 1
	resp := &response.FeedPost{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/feed/posts/link", params, asUserID, resp)

	require.NotNil(t, resp.Creator)
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, 1, resp.Creator.ID)
	assert.Nil(t, resp.Customer)
}

func (m *FeedTests) SubTestGetGoatFeedPosts(t *testing.T) {
	params := map[string]string{
		"goatUserID": "1",
	}
	asUserID := 2
	resp := &response.GetGoatFeedPosts{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/feed/posts/goat", params, asUserID, resp)

	require.NotNil(t, resp.FeedPosts)
	require.Equal(t, 2, len(resp.FeedPosts))
	assert.Equal(t, 2, resp.FeedPosts[0].ID)
	assert.Equal(t, 1, resp.FeedPosts[1].ID)
}

// Test that the creator user sees their own posts on home feed
func (m *FeedTests) SubTestGetMyHomeFeedPosts(t *testing.T) {
	params := map[string]string{}
	asUserID := 1
	resp := &response.GetUserFeedPosts{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/feed/posts/me", params, asUserID, resp)

	require.Equal(t, 2, len(resp.FeedPosts))
	assert.Equal(t, 2, resp.FeedPosts[0].ID)
	assert.Equal(t, 1, resp.FeedPosts[1].ID)
}

// Test that home feed is initially empty, follow creator, then
// see creator's posts on feed
func (m *FeedTests) SubTestFollowGoatThenGetHomeFeedPosts(t *testing.T) {
	params := map[string]string{}
	asUserID := 2

	shouldBeEmptyResp := &response.GetUserFeedPosts{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/feed/posts/me", params, asUserID, shouldBeEmptyResp)
	assert.Equal(t, 0, len(shouldBeEmptyResp.FeedPosts))

	followParams := map[string]interface{}{
		"userID": 1,
	}
	followResp := &response.User{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/follows/follow", followParams, &asUserID, followResp)

	resp := &response.GetUserFeedPosts{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/feed/posts/me", params, asUserID, resp)

	require.Equal(t, 2, len(resp.FeedPosts))

	assert.Equal(t, 2, resp.FeedPosts[0].ID)
	require.NotNil(t, resp.FeedPosts[0].Creator)
	assert.Equal(t, 1, resp.FeedPosts[0].Creator.ID)

	assert.Equal(t, 1, resp.FeedPosts[1].ID)
	require.NotNil(t, resp.FeedPosts[1].Creator)
	assert.Equal(t, 1, resp.FeedPosts[1].Creator.ID)
}

func (m *FeedTests) SubTestMakeFeedPost(t *testing.T) {
	textContent := "hihihihihihi feed post :) "
	params := map[string]interface{}{
		"textContent": textContent,
	}
	asUserID := 1
	resp := &response.FeedPost{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/feed/posts", params, &asUserID, resp)

	assert.Equal(t, 1, resp.Creator.ID)
	assert.Equal(t, textContent, resp.PostTextContent)
	assert.Equal(t, "", resp.Reaction)
	assert.Equal(t, 0, resp.NumComments)
}

func (m *FeedTests) SubTestDeleteFeedPost(t *testing.T) {
	postID := 1
	params := map[string]interface{}{
		"postID": postID,
	}

	asUserID := 1
	resp := map[string]interface{}{}
	test_util.MakeDeleteRequestAssert200(t, m.app, "/api/v1/feed/posts", params, asUserID, &resp)
}

func (m *FeedTests) SubTestUpvoteFeedPost(t *testing.T) {
	params := map[string]interface{}{
		"postID": 1,
	}

	asUserID := 1
	resp := response.Reaction{}
	test_util.MakePatchRequestAssert200(t, m.app, "/api/v1/feed/posts/upvote", params, &asUserID, &resp)

	assert.Equal(t, appconfig.Config.Vote.UpVote, resp.NewState)
	assert.Equal(t, int64(1), resp.NumUpvotes)
	assert.Equal(t, int64(0), resp.NumDownvotes)
}

func (m *FeedTests) SubTestDownvoteFeedPost(t *testing.T) {
	params := map[string]interface{}{
		"postID": 1,
	}

	asUserID := 1
	resp := response.Reaction{}
	test_util.MakePatchRequestAssert200(t, m.app, "/api/v1/feed/posts/downvote", params, &asUserID, &resp)

	assert.Equal(t, appconfig.Config.Vote.DownVote, resp.NewState)
	assert.Equal(t, int64(0), resp.NumUpvotes)
	assert.Equal(t, int64(1), resp.NumDownvotes)
}

// Tests that calling upvote twice returns the post to the
// original values
func (m *FeedTests) SubTestUpvoteFeedPostTwice(t *testing.T) {
	params := map[string]interface{}{
		"postID": 1,
	}

	asUserID := 1
	resp := response.Reaction{}
	test_util.MakePatchRequestAssert200(t, m.app, "/api/v1/feed/posts/upvote", params, &asUserID, &resp)
	resp = response.Reaction{}
	test_util.MakePatchRequestAssert200(t, m.app, "/api/v1/feed/posts/upvote", params, &asUserID, &resp)

	assert.Equal(t, appconfig.Config.Vote.Nil, resp.NewState)
	assert.Equal(t, int64(0), resp.NumUpvotes)
	assert.Equal(t, int64(0), resp.NumDownvotes)
}

// Tests that calling downvote twice returns the post to the
// original values
func (m *FeedTests) SubTestDownvoteFeedPostTwice(t *testing.T) {
	params := map[string]interface{}{
		"postID": 1,
	}

	asUserID := 1
	resp := response.Reaction{}
	test_util.MakePatchRequestAssert200(t, m.app, "/api/v1/feed/posts/downvote", params, &asUserID, &resp)
	resp = response.Reaction{}
	test_util.MakePatchRequestAssert200(t, m.app, "/api/v1/feed/posts/downvote", params, &asUserID, &resp)

	assert.Equal(t, appconfig.Config.Vote.Nil, resp.NewState)
	assert.Equal(t, int64(0), resp.NumUpvotes)
	assert.Equal(t, int64(0), resp.NumDownvotes)
}

// Test getting a post by ID that does not exist returns a 404.
func (m *FeedTests) SubTestGetFeedPostByIDNotFound(t *testing.T) {
	params := map[string]string{
		"postID": "100000",
	}

	asUserID := 1
	resp := &response.FeedPost{}
	test_util.MakeGetRequest(t, m.app, "/api/v1/feed/posts", 404, params, asUserID, resp)
}

func TestFeed(t *testing.T) {
	gtest.RunSubTests(t, &FeedTests{})
}
