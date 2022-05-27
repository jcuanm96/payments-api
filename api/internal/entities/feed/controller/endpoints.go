package controller

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	"github.com/VamaSingapore/vama-api/internal/vamawebsocket"
)

func MakeFeedPost(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)

	req := incomeRequest.(request.MakeFeedPost)

	res, err := svc.MakeFeedPost(ctx, req)

	return res, err
}

func MakeComment(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)

	req := incomeRequest.(request.MakeComment)

	res, err := svc.MakeComment(ctx, req)

	return res, err
}

func GetUserFeedPosts(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)

	req := incomeRequest.(request.GetUserFeedPosts)

	res, err := svc.GetUserFeedPosts(ctx, req)

	return res, err
}

func GetGoatFeedPosts(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)

	req := incomeRequest.(request.GetGoatFeedPosts)

	res, err := svc.GetGoatFeedPosts(ctx, req)

	return res, err
}

func UpvotePost(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)

	req := incomeRequest.(request.UpvotePost)

	res, err := svc.UpvotePost(ctx, req)

	return res, err
}

func DownvotePost(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)

	req := incomeRequest.(request.DownvotePost)

	res, err := svc.DownvotePost(ctx, req)

	return res, err
}

func GetFeedPostByID(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)

	req := incomeRequest.(request.GetFeedPostByID)

	res, err := svc.GetFeedPostByID(ctx, req)

	return res, err
}

func GetFeedPostByLinkSuffix(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)
	req := incomeRequest.(request.GetFeedPostByLinkSuffix)

	res, err := svc.GetFeedPostByLinkSuffix(ctx, req)

	return res, err
}

func PublicGetFeedPostByLinkSuffix(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)
	req := incomeRequest.(request.GetFeedPostByLinkSuffix)

	res, err := svc.PublicGetFeedPostByLinkSuffix(ctx, req)

	return res, err
}

func DeleteFeedPost(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)

	req := incomeRequest.(request.DeleteFeedPost)

	res, err := svc.DeleteFeedPost(ctx, req)

	return res, err
}

func GetFeedPostComments(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)

	req := incomeRequest.(request.GetFeedPostComments)

	res, err := svc.GetFeedPostComments(ctx, req)

	return res, err
}

func DeleteComment(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(feed.Usecase)

	req := incomeRequest.(request.DeleteComment)

	res, err := svc.DeleteComment(ctx, req)

	return res, err
}

func ReactWebsocket(uc interface{}, ctx context.Context, c *vamawebsocket.Conn, incomeRequest interface{}) error {
	svc := uc.(feed.Usecase)
	req := incomeRequest.(request.ReactWebsocket)

	err := svc.ReactWebsocket(ctx, c, req)

	return err
}
