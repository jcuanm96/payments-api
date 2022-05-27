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

const errDeletingComment = "Something went wrong when trying to delete comment."

func (svc *usecase) DeleteComment(ctx context.Context, req request.DeleteComment) (interface{}, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingComment,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errDeletingComment,
			"user was nil in DeleteComment",
		)
	}
	deleteErr := svc.repo.DeleteComment(ctx, req.CommentID, user.ID)
	if deleteErr == constants.ErrNoRowsAffected {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You cannot delete this comment.",
			fmt.Sprintf("Could not delete user %d's comment with ID: %d: %v", user.ID, req.CommentID, deleteErr),
		)
	} else if deleteErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingComment,
			fmt.Sprintf("Could not delete user %d's comment with ID: %d: %v", user.ID, req.CommentID, deleteErr),
		)
	}
	res := response.DeleteFeedPost{}
	return &res, nil
}
