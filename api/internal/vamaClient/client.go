package vama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL     string
	accessToken string
	httpClient  *http.Client
}

func NewClient(baseURL string) *Client {
	c := Client{
		BaseURL: baseURL,
	}
	c.httpClient = &http.Client{
		Timeout: time.Second * 5,
	}
	return &c
}

func (c *Client) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

func (c *Client) prepareHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
}

func (c *Client) PrepareUrl(pathEncodedUrl string) *url.URL {
	urlVal := &url.URL{
		Scheme:  "http",
		Host:    c.BaseURL,
		Path:    "/api/v1" + pathEncodedUrl,
		RawPath: "/api/v1" + pathEncodedUrl,
	}

	return urlVal
}

func (c *Client) get(config *url.URL, rawQuery string, resp interface{}) error {
	req, reqErr := http.NewRequest("GET", config.String(), nil)
	if reqErr != nil {
		return reqErr
	}

	c.prepareHeader(req)
	req.URL.RawQuery = rawQuery

	httpResp, getHttpErr := c.httpClient.Do(req)
	if getHttpErr != nil {
		return getHttpErr
	}

	defer httpResp.Body.Close()

	processedErr := CheckVamaError(httpResp)
	if processedErr != nil {
		return processedErr
	}

	return json.NewDecoder(httpResp.Body).Decode(resp)
}

func (c *Client) post(config *url.URL, apiReq interface{}, resp interface{}) error {
	body, marshalErr := json.Marshal(apiReq)
	if marshalErr != nil {
		return marshalErr
	}

	req, reqErr := http.NewRequest("POST", config.String(), bytes.NewBuffer(body))
	if reqErr != nil {
		return reqErr
	}

	c.prepareHeader(req)

	req.URL.RawQuery = url.Values{}.Encode()

	httpResp, postHttpErr := c.httpClient.Do(req)
	if postHttpErr != nil {
		return postHttpErr
	}

	defer httpResp.Body.Close()
	processedErr := CheckVamaError(httpResp)
	if processedErr != nil {
		return processedErr
	}

	return json.NewDecoder(httpResp.Body).Decode(resp)
}

type VamaError struct {
	Code    int    `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Implement error interface
func (s VamaError) Error() string {
	if s.Code != 200 && s.Code != 0 {
		return fmt.Sprintf("Vama: %d - %s", s.Code, s.Message)
	}
	return "{}"
}

func CheckVamaError(httpResp *http.Response) error {
	if httpResp.StatusCode != 200 {
		errorMessageBody := VamaError{}
		err := json.NewDecoder(httpResp.Body).Decode(&errorMessageBody)
		if err != nil {
			return fmt.Errorf("vama client error: %s", err)
		}

		return errorMessageBody
	}
	return nil
}
