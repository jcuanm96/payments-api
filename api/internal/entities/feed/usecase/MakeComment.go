package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errPostingComment = "Something went wrong when posting comment."

func (svc *usecase) MakeComment(ctx context.Context, req request.MakeComment) (*response.Comment, error) {
	user, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errPostingComment,
			fmt.Sprintf("Could not find user in the current context. Err: %v", getCurrUserErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errPostingComment,
			"user was nil in MakeComment",
		)
	}

	isBlocked, isBlockedErr := svc.repo.IsUserBlockedByPoster(ctx, user.ID, req.PostID)
	if isBlockedErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errPostingComment,
			fmt.Sprintf("Error checking if user %d is blocked: %v", user.ID, isBlockedErr),
		)
	}

	if isBlocked {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You cannot comment on this post.",
			"user is blocked by creator of post",
		)
	}

	commentMetadata, makeCommentErr := svc.repo.MakeComment(ctx, req, user.ID)
	if makeCommentErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errPostingComment,
			fmt.Sprintf("Could not make comment for user %d. Err: %v", user.ID, makeCommentErr),
		)
	} else if commentMetadata == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errPostingComment,
			"commentMetadata was nil in MakeComment",
		)
	}

	res := response.Comment{
		CommentMetadata: *commentMetadata,
		User:            *user,
	}
	return &res, nil
}
