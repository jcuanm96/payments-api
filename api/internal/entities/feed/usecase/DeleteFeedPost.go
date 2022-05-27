package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errDeletingFeedPost = "Something went wrong when trying to delete comment."

func (svc *usecase) DeleteFeedPost(ctx context.Context, req request.DeleteFeedPost) (*response.DeleteFeedPost, error) {
	user, getCurrentUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrentUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingFeedPost,
			fmt.Sprintf("Could not find user in the current context: %v", getCurrentUserErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errDeletingFeedPost,
			"user was nil in DeleteFeedPost",
		)
	}
	deletePostErr := svc.repo.DeleteFeedPost(ctx, req.PostID, user.ID)
	if deletePostErr == constants.ErrNoRowsAffected {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You cannot delete this post.",
			fmt.Sprintf("User %d cannot delete post %d.", user.ID, req.PostID),
		)
	} else if deletePostErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingFeedPost,
			fmt.Sprintf("Could not delete user %d's post with ID: %d: %v", user.ID, req.PostID, deletePostErr),
		)
	}
	res := response.DeleteFeedPost{}
	return &res, nil
}
