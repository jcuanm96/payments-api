package follow

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

type Usecase interface {
	Follow(ctx context.Context, goatUserID int) (*response.User, error)
	Unfollow(ctx context.Context, userToUnfollowID int) error
	IsFollowing(ctx context.Context, goatUserID int) (*response.IsFollowing, error)
	GetFollowedGoats(ctx context.Context, req request.GetFollowedGoats) (*response.GetFollowedGoats, error)
}
