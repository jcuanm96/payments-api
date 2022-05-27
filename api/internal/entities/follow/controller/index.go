package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {

	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/follows/follow",
			Method:         "post",
			Handler:        Follow,
			RequestDecoder: decodeFollow,
		},
		{
			Path:           "/follows/unfollow",
			Method:         "post",
			Handler:        Unfollow,
			RequestDecoder: decodeUnfollow,
		},
		{
			Path:           "/follows/check",
			Method:         "get",
			Handler:        IsFollowing,
			RequestDecoder: decodeIsFollowing,
		},
		{
			Path:           "/follows/goats",
			Method:         "get",
			Handler:        GetFollowedGoats,
			RequestDecoder: decodeGetFollowedGoats,
		},
	}

	return res
}
