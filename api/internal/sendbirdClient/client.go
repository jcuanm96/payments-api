package sendbird

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

type Client interface {
	GetUser(userID int) (*User, error)
	CreateUser(createUserReq *CreateUserParams) (*User, error)
	DeleteUser(userID int) error
	UpsertUserMetadata(userID int, newMetadata interface{}) error
	UpdateUser(userID int, req *UpdateUserParams) (*User, error)
	BlockUser(userID, blockUserID int) (*User, error)
	UnblockUser(userID, unblockUserID int) error

	ListGroupChannelMessages(req ListGroupChannelMessagesParams) (*ListMessages, error)
	ViewGroupChannelMessage(channelURL string, messageID string) (*Message, error)
	SendMessage(channelURL string, body *SendMessageParams) (*Message, error)
	SearchMessages(req SearchMessagesParams) (*SearchMessagesResponse, error)

	GetGroupChannel(channelURL string, req GetGroupChannelParams) (*GroupChannel, error)
	ListGroupChannels(req ListGroupChannelsParams) ([]GroupChannel, error)
	ListMyGroupChannels(userID int, req ListMyGroupChannelsParams) ([]GroupChannel, error)
	CreateGroupChannel(ctx context.Context, req *CreateGroupChannelParams) (*GroupChannel, error)
	UpdateGroupChannel(channelURL string, req *UpdateGroupChannelParams, shouldUpsert bool) (*GroupChannel, error)
	IsMemberInGroupChannel(channelURL string, userID string) (*IsMemberInGroupChannel, error)
	InviteGroupChannelMembers(channelURL string, req *InviteGroupChannelMembers) (*GroupChannel, error)
	ListGroupChannelMembers(channelURL string, req ListGroupChannelMembersParams) (*ListGroupChannelMembersResponse, error)
	DeleteGroupChannel(channelURL string) error
	JoinGroupChannel(channelURL string, req *JoinGroupChannelParams) error
	LeaveGroupChannel(channelURL string, req *LeaveGroupChannelParams) error
	BanUserFromGroupChannel(channelURL string, req *BanUserFromGroupChannelParams) error
	UnbanUserFromGroupChannel(channelURL string, userID string) error
}

type client struct {
	httpClient    *http.Client
	apiKey        string
	applicationID string
}

func NewClient() Client {
	c := client{}
	c.httpClient = &http.Client{
		Timeout: time.Second * 5,
	}
	c.apiKey = appconfig.Config.Sendbird.MasterAPIKey
	c.applicationID = appconfig.Config.Sendbird.ApplicationID

	return &c
}

func (c *client) prepareHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Token", c.apiKey)
}

func (c *client) PrepareUrl(pathEncodedUrl string) *url.URL {
	urlVal := &url.URL{
		Scheme:  "https",
		Host:    fmt.Sprintf("api-%s.sendbird.com", c.applicationID),
		Path:    "/v3" + pathEncodedUrl,
		RawPath: "/v3" + pathEncodedUrl,
	}
	return urlVal
}

func (c *client) get(config *url.URL, rawQuery string, resp interface{}) error {
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

	processedErr := CheckSendbirdError(httpResp)
	if processedErr != nil {
		return processedErr
	}

	return json.NewDecoder(httpResp.Body).Decode(resp)
}

func (c *client) post(config *url.URL, apiReq interface{}, resp interface{}) error {
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
	processedErr := CheckSendbirdError(httpResp)
	if processedErr != nil {
		return processedErr
	}

	return json.NewDecoder(httpResp.Body).Decode(resp)
}

func (c *client) delete(config *url.URL, rawQueryString string, resp interface{}) error {
	req, reqErr := http.NewRequest("DELETE", config.String(), nil)
	if reqErr != nil {
		return reqErr
	}

	c.prepareHeader(req)
	req.URL.RawQuery = rawQueryString

	httpResp, deleteErr := c.httpClient.Do(req)
	if deleteErr != nil {
		return deleteErr
	}

	defer httpResp.Body.Close()

	processedErr := CheckSendbirdError(httpResp)
	if processedErr != nil {
		return processedErr
	}

	return json.NewDecoder(httpResp.Body).Decode(resp)
}

func (c *client) put(config *url.URL, apiReq interface{}, resp interface{}) error {
	body, marshalErr := json.Marshal(apiReq)
	if marshalErr != nil {
		return marshalErr
	}
	req, reqErr := http.NewRequest("PUT", config.String(), bytes.NewBuffer(body))
	if reqErr != nil {
		return reqErr
	}

	c.prepareHeader(req)

	req.URL.RawQuery = url.Values{}.Encode()

	httpResp, putHttpErr := c.httpClient.Do(req)
	if putHttpErr != nil {
		return putHttpErr
	}

	defer httpResp.Body.Close()

	processedErr := CheckSendbirdError(httpResp)
	if processedErr != nil {
		return processedErr
	}

	return json.NewDecoder(httpResp.Body).Decode(resp)
}

func (c *client) createGroupChannelWithFile(ctx context.Context, req *CreateGroupChannelParams) (*GroupChannel, error) {
	if req.CoverFile == nil {
		return nil, errors.New("cover file not specified")
	}

	coverFile := req.CoverFile
	userIDs := req.UserIDs
	operatorIDs := req.OperatorIDs
	// set these fields to nil to not include them in the marshal/unmarshal of the params
	// (works because of the `omitempty` tag)
	req.CoverFile = nil
	req.UserIDs = nil
	req.OperatorIDs = nil

	reqBytes, marshalErr := json.Marshal(req)
	if marshalErr != nil {
		return nil, marshalErr
	}
	paramsMap := map[string]interface{}{}
	unmarshalErr := json.Unmarshal(reqBytes, &paramsMap)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	file, fileErr := coverFile.Open()
	if fileErr != nil {
		return nil, fileErr
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer writer.Close()
	part, createFormFileErr := writer.CreateFormFile("cover_file", coverFile.Filename)
	if createFormFileErr != nil {
		return nil, createFormFileErr
	}
	io.Copy(part, file)

	// Add all non-complex params
	for key, val := range paramsMap {
		writeFieldErr := writer.WriteField(key, fmt.Sprintf("%v", val))
		if writeFieldErr != nil {
			vlog.Errorf(ctx, "Error adding field key %s value %s: %v", key, val, writeFieldErr)
		}
	}
	// Add array params
	for _, userID := range userIDs {
		writer.WriteField("user_ids", fmt.Sprint(userID))
	}

	for _, operatorID := range operatorIDs {
		writer.WriteField("operator_ids", fmt.Sprint(operatorID))
	}

	pathString := "/group_channels"
	parsedURL := c.PrepareUrl(pathString)
	httpReq, httpReqErr := http.NewRequest("POST", parsedURL.String(), body)
	if httpReqErr != nil {
		return nil, httpReqErr
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("Api-Token", c.apiKey)

	resp := &GroupChannel{}
	httpReq.URL.RawQuery = url.Values{}.Encode()

	httpResp, postHttpErr := c.httpClient.Do(httpReq)
	if postHttpErr != nil {
		return nil, postHttpErr
	}

	defer httpResp.Body.Close()
	processedErr := CheckSendbirdError(httpResp)
	if processedErr != nil {
		return nil, processedErr
	}

	decodeErr := json.NewDecoder(httpResp.Body).Decode(resp)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return resp, nil
}
