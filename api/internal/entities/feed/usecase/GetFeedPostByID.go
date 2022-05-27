package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errRetrievingPost = "Something went wrong when trying to retrieve feed post."

func (svc *usecase) GetFeedPostByID(ctx context.Context, req request.GetFeedPostByID) (*response.FeedPost, error) {
	user, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errRetrievingPost,
			fmt.Sprintf("Could not find user in the current context. Err: %v", getUserErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errRetrievingPost,
			"user was nil in GetFeedPostByID",
		)
	}

	params := feed.GetFeedPostsParams{
		IsTextContentFullLength: true,
		PostID:                  &req.PostID,
	}
	res, sqlErr := svc.repo.GetFeedPosts(ctx, user.ID, params)

	if sqlErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errRetrievingPost,
			fmt.Sprintf("Could not get feed post %d from db: %v", req.PostID, sqlErr),
		)
	}

	if len(res) == 0 || res[0].ID <= 0 {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"Sorry, that post doesn't seem to exist.",
			"No feed post returned or ID was <= 0",
		)
	}

	if len(res) > 1 {
		msg := fmt.Sprintf("GetFeedPostByID did not return exactly 1 valid post. db response: %v", res)
		telegram.TelegramClient.SendMessage(msg)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errRetrievingPost,
			msg,
		)
	}

	svc.addPreviewMessagesToPost(ctx, &res[0])

	return &res[0], nil
}
