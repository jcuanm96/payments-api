package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {
	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/feed/posts",
			Method:         "post",
			Handler:        MakeFeedPost,
			RequestDecoder: decodeMakeFeedPost,
		},
		{
			Path:           "/feed/posts/comments",
			Method:         "post",
			Handler:        MakeComment,
			RequestDecoder: decodeMakeComment,
		},
		{
			Path:           "/feed/posts/upvote",
			Method:         "patch",
			Handler:        UpvotePost,
			RequestDecoder: decodeUpvotePost,
		},
		{
			Path:           "/feed/posts/downvote",
			Method:         "patch",
			Handler:        DownvotePost,
			RequestDecoder: decodeDownvotePost,
		},
		{
			Path:           "/feed/posts/me",
			Method:         "get",
			Handler:        GetUserFeedPosts,
			RequestDecoder: decodeGetUserFeedPosts,
		},
		{
			Path:           "/feed/posts/goat",
			Method:         "get",
			Handler:        GetGoatFeedPosts,
			RequestDecoder: decodeGetGoatFeedPosts,
		},
		{
			Path:           "/feed/posts",
			Method:         "get",
			Handler:        GetFeedPostByID,
			RequestDecoder: decodeGetFeedPostByID,
		},
		{
			Path:           "/feed/posts/link",
			Method:         "get",
			Handler:        GetFeedPostByLinkSuffix,
			RequestDecoder: decodeGetFeedPostByLinkSuffix,
		},
		{
			Path:           "/feed/posts/link",
			Method:         "get",
			Handler:        PublicGetFeedPostByLinkSuffix,
			Version:        constants.PUBLIC_V1,
			RequestDecoder: decodeGetFeedPostByLinkSuffix,
		},
		{
			Path:           "/feed/posts",
			Method:         "delete",
			Handler:        DeleteFeedPost,
			RequestDecoder: decodeDeleteFeedPost,
		},
		{
			Path:           "/feed/posts/comments",
			Method:         "get",
			Handler:        GetFeedPostComments,
			RequestDecoder: decodeGetFeedPostComments,
		},
		{
			Path:           "/feed/comments",
			Method:         "delete",
			Handler:        DeleteComment,
			RequestDecoder: decodeDeleteComment,
		},
		{
			Path:             "/feed/reacts",
			Method:           "get",
			WebsocketHandler: ReactWebsocket,
			RequestDecoder:   decodeReactWebsocket,
			Version:          constants.WEBSOCKET_V1,
		},
	}

	return res
}
