package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errGettingFeedPosts = "Something went wrong when trying to retrieve feed posts."

func (svc *usecase) GetUserFeedPosts(ctx context.Context, req request.GetUserFeedPosts) (*response.GetUserFeedPosts, error) {
	user, getCurrentUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrentUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errGettingFeedPosts,
			fmt.Sprintf("Could not find user in the current context. Err: %v", getCurrentUserErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errGettingFeedPosts,
			"user was nil in GetUserFeedPosts",
		)
	}

	params := feed.GetFeedPostsParams{
		CursorPostID: &req.CursorPostID,
		Limit:        &req.Limit,
	}
	userFeedPosts, getFeedPostsErr := svc.repo.GetFeedPosts(ctx, user.ID, params)
	if getFeedPostsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errGettingFeedPosts,
			fmt.Sprintf("Could not get user home feed posts for user %d: %v", user.ID, getFeedPostsErr),
		)
	}

	svc.addPreviewMessagesToPosts(ctx, userFeedPosts)

	res := &response.GetUserFeedPosts{
		FeedPosts: userFeedPosts,
	}

	return res, nil
}

func (svc *usecase) addPreviewMessagesToPost(ctx context.Context, post *response.FeedPost) error {
	if post.Conversation.ChannelID == "" {
		return nil
	}
	listChatMessagesReq := request.ListChatMessages{
		SendBirdChannelID: post.Conversation.ChannelID,
		MessageTsFrom:     int64(post.Conversation.StartTS),
		Limit:             200,
	}

	if post.Conversation.EndTS != 0 {
		endTS := int64(post.Conversation.EndTS)
		listChatMessagesReq.MessageTsTo = &endTS
	}

	listChatMessagesRes, listMessagesHttpErr := (*svc.chat).ListChatMessages(ctx, listChatMessagesReq)
	if listMessagesHttpErr != nil {
		vlog.Errorf(ctx, "Error getting Sendbird messages for preview for post %d. Err: %v", post.ID, listMessagesHttpErr)
		return listMessagesHttpErr
	}

	post.Conversation.PreviewMessages = listChatMessagesRes.Messages
	return nil
}
func (svc *usecase) addPreviewMessagesToPosts(ctx context.Context, posts []response.FeedPost) {
	var wg sync.WaitGroup
	wg.Add(len(posts))
	for i := range posts {
		go func(post *response.FeedPost) {
			defer wg.Done()
			svc.addPreviewMessagesToPost(ctx, post)
		}(&posts[i])
	}

	wg.Wait()
}
