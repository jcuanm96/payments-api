package follow

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
)

type Repository interface {
	baserepo.BaseRepository
	Follow(ctx context.Context, userID int, goatUserID int) error
	Unfollow(ctx context.Context, userID int, userToUnfollowID int) error
	IsFollowing(ctx context.Context, userID, goatUserID int) (bool, error)
	GetFollowedGoats(ctx context.Context, userID int, req request.GetFollowedGoats) ([]response.User, error)
}
