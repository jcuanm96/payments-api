package controller

import (
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func decodeMakeFeedPost(c *fiber.Ctx) (interface{}, error) {
	var p request.MakeFeedPost

	c.BodyParser(&p)

	const maxPostTextContentLength = 40000
	if len(p.TextContent) > maxPostTextContentLength {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Post too long, please keep it under %d characters", maxPostTextContentLength),
			fmt.Sprintf("Post too long, please keep it under %d characters", maxPostTextContentLength),
		)
	}

	file, imageErr := c.FormFile("image")
	if imageErr != nil && p.Image != nil {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong uploading your image.",
			fmt.Sprintf("Error getting image: %v", imageErr),
		)
	}

	p.Image = file

	return p, nil
}

func decodeGetUserFeedPosts(c *fiber.Ctx) (interface{}, error) {
	var q request.GetUserFeedPosts

	c.QueryParser(&q)

	if q.Limit == 0 {
		q.Limit = 10
	}

	return q, nil
}

func decodeGetGoatFeedPosts(c *fiber.Ctx) (interface{}, error) {
	var q request.GetGoatFeedPosts

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.GoatUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"goatUserID is required.",
		)
	}

	if q.Limit == 0 {
		q.Limit = 10
	}

	return q, nil
}

func decodeGetFeedPostByID(c *fiber.Ctx) (interface{}, error) {
	var q request.GetFeedPostByID

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.PostID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"postID is required.",
		)
	}

	return q, nil
}

func decodeGetFeedPostByLinkSuffix(c *fiber.Ctx) (interface{}, error) {
	var q request.GetFeedPostByLinkSuffix

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.LinkSuffix) {
		return nil, httperr.New(
			404,
			http.StatusNotFound,
			"Sorry, that post doesn't seem to exist.",
			"linkSuffix is required.",
		)
	}

	return q, nil
}

func decodeGetFeedPostComments(c *fiber.Ctx) (interface{}, error) {
	var q request.GetFeedPostComments
	c.QueryParser(&q)

	if utils.IsEmptyValue(q.PostID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"postID is required.",
		)
	}

	return q, nil
}

func decodeMakeComment(c *fiber.Ctx) (interface{}, error) {
	var p request.MakeComment

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.PostID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"postID is required.",
		)
	}

	if utils.IsEmptyValue(p.Text) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"text is required.",
		)
	}

	const maxPostCommentLength = 10000
	if len(p.Text) > maxPostCommentLength {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Comment too long, please keep it under %d characters", maxPostCommentLength),
			fmt.Sprintf("Comment too long, please keep it under %d characters", maxPostCommentLength),
		)
	}

	return p, nil
}

func decodeDownvotePost(c *fiber.Ctx) (interface{}, error) {
	var p request.DownvotePost

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.PostID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"postID is required.",
		)
	}

	return p, nil
}

func decodeUpvotePost(c *fiber.Ctx) (interface{}, error) {
	var p request.UpvotePost

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.PostID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"postID is required.",
		)
	}

	return p, nil
}

func decodeDeleteFeedPost(c *fiber.Ctx) (interface{}, error) {
	var p request.DeleteFeedPost

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.PostID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"postID is required.",
		)
	}

	return p, nil
}

func decodeDeleteComment(c *fiber.Ctx) (interface{}, error) {
	var p request.DeleteComment

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.CommentID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"commentID is required.",
		)
	}

	return p, nil
}

func decodeReactWebsocket(c *fiber.Ctx) (interface{}, error) {
	var p request.ReactWebsocket

	c.QueryParser(&p)
	return p, nil
}
