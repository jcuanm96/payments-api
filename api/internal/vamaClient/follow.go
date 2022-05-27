package vama

import (
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (c *Client) Follow(req request.Follow) (*response.User, error) {
	pathString := "/follows/follow"
	parsedURL := c.PrepareUrl(pathString)

	result := response.User{}
	postReqErr := c.post(parsedURL, req, &result)
	if postReqErr != nil {
		return nil, postReqErr
	}

	return &result, nil
}
