package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	"github.com/VamaSingapore/vama-api/internal/entities/sharing"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) PublicGetFeedPostByLinkSuffix(ctx context.Context, req request.GetFeedPostByLinkSuffix) (*response.PublicFeedPost, error) {
	params := feed.GetFeedPostsParams{
		IsTextContentFullLength: true,
		LinkSuffix:              &req.LinkSuffix,
	}

	emptyUserID := 0
	posts, sqlErr := svc.repo.GetFeedPosts(ctx, emptyUserID, params)

	if sqlErr != nil {
		vlog.Errorf(ctx, "Error getting feed post %s from db: %v", req.LinkSuffix, sqlErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errRetrievingPost,
		)
	}

	if len(posts) == 0 || posts[0].ID <= 0 {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"Sorry, that post doesn't seem to exist.",
			"No feed post returned or ID was <= 0",
		)
	}

	if len(posts) > 1 {
		msg := fmt.Sprintf("PublicGetFeedPostByLinkSuffix did not return exactly 1 valid post. db response: %v", posts)
		vlog.Errorf(ctx, msg)
		telegram.TelegramClient.SendMessage(msg)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errRetrievingPost,
		)
	}

	svc.addPreviewMessagesToPost(ctx, &posts[0])
	post := posts[0]

	deepLinkBaseURL := fmt.Sprintf("https://%s", appconfig.Config.Gcloud.RedirectBaseURL)
	deepLink := fmt.Sprintf(constants.FEED_POST_BASE_URL_F, deepLinkBaseURL, req.LinkSuffix)
	dynamicLink, getDynamicLinkErr := sharing.CreateDynamicLink(ctx, deepLink)
	if getDynamicLinkErr != nil {
		vlog.Errorf(ctx, "Error creating dynamic link for feed post suffix %s: %v", req.LinkSuffix, getDynamicLinkErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	}

	publicCreator := response.User{
		FirstName:     post.Creator.FirstName,
		LastName:      post.Creator.LastName,
		Username:      post.Creator.Username,
		ProfileAvatar: post.Creator.ProfileAvatar,
	}

	res := &response.PublicFeedPost{
		PostCreatedAt:    post.PostCreatedAt,
		PostNumUpvotes:   post.PostNumUpvotes,
		PostNumDownvotes: post.PostNumDownvotes,
		PostTextContent:  post.PostTextContent,
		NumComments:      post.NumComments,
		IsConversation:   post.Conversation.ChannelID != "",
		Image:            post.Image,
		Link:             post.Link,
		Creator:          publicCreator,
		DynamicLink:      dynamicLink,
	}

	if post.Customer != nil {
		publicCustomer := &response.User{
			FirstName:     post.Customer.FirstName,
			LastName:      post.Customer.LastName,
			Username:      post.Customer.Username,
			ProfileAvatar: post.Customer.ProfileAvatar,
		}
		res.Customer = publicCustomer
	}

	return res, nil
}
