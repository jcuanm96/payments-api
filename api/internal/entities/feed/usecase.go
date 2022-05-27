package feed

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/internal/vamawebsocket"
)

type Usecase interface {
	MakeFeedPost(ctx context.Context, req request.MakeFeedPost) (*response.FeedPost, error)
	MakeComment(ctx context.Context, req request.MakeComment) (*response.Comment, error)
	UpvotePost(ctx context.Context, req request.UpvotePost) (interface{}, error)
	DownvotePost(ctx context.Context, req request.DownvotePost) (interface{}, error)
	GetUserFeedPosts(ctx context.Context, req request.GetUserFeedPosts) (*response.GetUserFeedPosts, error)
	GetGoatFeedPosts(ctx context.Context, req request.GetGoatFeedPosts) (*response.GetGoatFeedPosts, error)
	GetFeedPostByID(ctx context.Context, req request.GetFeedPostByID) (*response.FeedPost, error)
	GetFeedPostByLinkSuffix(ctx context.Context, req request.GetFeedPostByLinkSuffix) (*response.FeedPost, error)
	PublicGetFeedPostByLinkSuffix(ctx context.Context, req request.GetFeedPostByLinkSuffix) (*response.PublicFeedPost, error)
	DeleteFeedPost(ctx context.Context, req request.DeleteFeedPost) (*response.DeleteFeedPost, error)
	GetFeedPostComments(ctx context.Context, req request.GetFeedPostComments) (*response.GetFeedPostComments, error)
	DeleteComment(ctx context.Context, req request.DeleteComment) (interface{}, error)
	GenerateReactionStruct(ctx context.Context, reactionType string, userID int, postID int, runnable utils.Runnable) (*Reaction, error)
	ReactWebsocket(ctx context.Context, conn *vamawebsocket.Conn, req request.ReactWebsocket) error
}
