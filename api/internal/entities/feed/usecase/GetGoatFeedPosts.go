package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetGoatFeedPosts(ctx context.Context, req request.GetGoatFeedPosts) (*response.GetGoatFeedPosts, error) {
	user, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errGettingFeedPosts,
			fmt.Sprintf("Could not find user in the current context. Err: %v", getUserErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errGettingFeedPosts,
			"user was nil in GetGoatFeedPosts",
		)
	}

	params := feed.GetFeedPostsParams{
		GoatUserID:   &req.GoatUserID,
		CursorPostID: &req.CursorPostID,
		Limit:        &req.Limit,
	}
	goatFeedPosts, sqlErr := svc.repo.GetFeedPosts(ctx, user.ID, params)
	if sqlErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errGettingFeedPosts,
			fmt.Sprintf("Error getting creator feed posts in db for creator %d: %v", req.GoatUserID, sqlErr),
		)
	}

	svc.addPreviewMessagesToPosts(ctx, goatFeedPosts)

	res := &response.GetGoatFeedPosts{
		FeedPosts: goatFeedPosts,
	}

	return res, nil
}
