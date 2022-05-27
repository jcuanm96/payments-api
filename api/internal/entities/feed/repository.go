package feed

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	MasterNode() *pgxpool.Pool
	baserepo.BaseRepository
	MakeComment(ctx context.Context, req request.MakeComment, userID int) (*response.CommentMetadata, error)
	GetReaction(ctx context.Context, userID int, postID int, runnable utils.Runnable) (*string, error)
	React(ctx context.Context, reaction Reaction, tx pgx.Tx) (*ReactionMetadata, error)

	GenerateFeedPostLinkSuffix(ctx context.Context) (string, error)
	UpsertFeedPost(ctx context.Context, req request.MakeFeedPost, userID int, image response.PostImage, linkSuffix string) (int, error)
	GetFeedPosts(ctx context.Context, userID int, params GetFeedPostsParams) ([]response.FeedPost, error)
	DeleteFeedPost(ctx context.Context, postID int, userID int) error

	GetFeedPostComments(ctx context.Context, req request.GetFeedPostComments, runnable utils.Runnable) ([]response.Comment, error)
	DeleteComment(ctx context.Context, commentID int, userID int) error

	IsUserBlockedByPoster(ctx context.Context, userID int, postID int) (bool, error)
}
