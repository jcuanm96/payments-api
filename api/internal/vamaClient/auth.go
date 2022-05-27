package vama

import (
	"net/url"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (c *Client) SignUpSMS(req request.SignupSMS) (*response.AuthSuccess, error) {
	pathString := "/auth/v1/sign-up/sms"
	urlVal := &url.URL{
		Scheme:  "http",
		Host:    c.BaseURL,
		Path:    pathString,
		RawPath: pathString,
	}

	result := response.AuthSuccess{}
	postReqErr := c.post(urlVal, req, &result)
	if postReqErr != nil {
		return nil, postReqErr
	}

	return &result, nil
}

func (c *Client) SignInSMS(req request.SignInSMS) (*response.AuthSuccess, error) {
	pathString := "/auth/v1/sign-in/sms"
	urlVal := &url.URL{
		Scheme:  "http",
		Host:    c.BaseURL,
		Path:    pathString,
		RawPath: pathString,
	}

	result := response.AuthSuccess{}
	postReqErr := c.post(urlVal, req, &result)
	if postReqErr != nil {
		return nil, postReqErr
	}

	return &result, nil
}
