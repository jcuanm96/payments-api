package vama

import (
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (c *Client) GetHomeFeedPosts(req request.GetUserFeedPosts) (*response.GetUserFeedPosts, error) {
	pathString := "/feed/posts/me"
	parsedURL := c.PrepareUrl(pathString)

	queryValues, queryValuesErr := ParseValues(req)
	if queryValuesErr != nil {
		return nil, queryValuesErr
	}
	rawQuery := queryValues.Encode()

	result := response.GetUserFeedPosts{}
	postReqErr := c.get(parsedURL, rawQuery, &result)
	if postReqErr != nil {
		return nil, postReqErr
	}

	return &result, nil
}

func (c *Client) GetFeedPostByID(postID int) (*response.FeedPost, error) {
	pathString := "/feed/posts"
	parsedURL := c.PrepareUrl(pathString)

	queryString := fmt.Sprintf("postID=%d", postID)

	result := response.FeedPost{}
	postReqErr := c.get(parsedURL, queryString, &result)
	if postReqErr != nil {
		return nil, postReqErr
	}

	return &result, nil
}

func (c *Client) MakeComment(req request.MakeComment) error {
	pathString := "/feed/posts/comments"
	parsedURL := c.PrepareUrl(pathString)

	result := struct{}{}
	postReqErr := c.post(parsedURL, req, &result)
	if postReqErr != nil {
		return postReqErr
	}

	return nil
}
