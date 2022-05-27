package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetFeedPostComments(ctx context.Context, req request.GetFeedPostComments) (*response.GetFeedPostComments, error) {
	comments, getCommentsErr := svc.repo.GetFeedPostComments(ctx, req, svc.repo.MasterNode())
	if getCommentsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to retrieve comments.",
			fmt.Sprintf("Could not get feed comments for post %d: %v", req.PostID, getCommentsErr),
		)
	}

	res := &response.GetFeedPostComments{
		Comments: comments,
	}

	return res, nil
}
