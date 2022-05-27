package sendbird

import (
	"fmt"
)

type User struct {
	UserID      string      `json:"user_id"`
	NickName    string      `json:"nickname"`
	ProfileURL  string      `json:"profile_url"`
	ProfileFile []byte      `json:"profile_file"`
	AccessToken string      `json:"access_token"`
	IsActive    bool        `json:"is_active"`
	IsOnline    bool        `json:"is_online"`
	LastSeenAt  int64       `json:"last_seen_at"`
	Metadata    interface{} `json:"metadata"`
	State       string      `json:"state,omitempty"` // Either 'joined' or 'invited' in webhook events
}

type CreateUserParams struct {
	UserID           int    `json:"user_id"`
	Nickname         string `json:"nickname"`
	ProfileURL       string `json:"profile_url"`
	IssueAccessToken bool   `json:"issue_access_token"`
}

func (c *client) DeleteUser(userID int) error {
	pathString := fmt.Sprintf("/users/%d", userID)
	parsedURL := c.PrepareUrl(pathString)

	result := SendbirdErrorResponse{}
	deleteReqErr := c.delete(parsedURL, "", &result)
	if deleteReqErr != nil {
		return deleteReqErr
	}

	return nil
}

func (c *client) CreateUser(createUserReq *CreateUserParams) (*User, error) {
	pathString := "/users"
	parsedURL := c.PrepareUrl(pathString)

	result := User{}
	postReqErr := c.post(parsedURL, createUserReq, &result)
	if postReqErr != nil {
		return nil, postReqErr
	}

	return &result, nil
}

// newMetadata should be a struct with json annotations that can be marshaled
func (c *client) UpsertUserMetadata(userID int, newMetadata interface{}) error {
	pathString := fmt.Sprintf("/users/%d/metadata", userID)
	parsedURL := c.PrepareUrl(pathString)

	req := map[string]interface{}{
		"metadata": newMetadata,
		"upsert":   true,
	}

	resp := map[string]string{}
	putReqErr := c.put(parsedURL, req, &resp)

	return putReqErr
}

type UpdateUserParams struct {
	Nickname   string `json:"nickname,omitempty"`
	ProfileURL string `json:"profile_url,omitempty"`
}

func (c *client) UpdateUser(userID int, req *UpdateUserParams) (*User, error) {
	pathString := fmt.Sprintf("/users/%d", userID)
	parsedURL := c.PrepareUrl(pathString)

	resp := &User{}
	putReqErr := c.put(parsedURL, req, resp)

	return resp, putReqErr
}

func (c *client) GetUser(userID int) (*User, error) {
	pathString := fmt.Sprintf("/users/%d", userID)
	parsedURL := c.PrepareUrl(pathString)

	result := User{}
	getReqErr := c.get(parsedURL, "", &result)
	if getReqErr != nil {
		return nil, getReqErr
	}

	return &result, nil
}

type BlockUserParams struct {
	TargetID int `json:"target_id"`
}

func (c *client) BlockUser(userID, blockUserID int) (*User, error) {
	pathString := fmt.Sprintf("/users/%d/block", userID)
	parsedURL := c.PrepareUrl(pathString)

	req := BlockUserParams{TargetID: blockUserID}

	result := User{}
	getReqErr := c.post(parsedURL, req, &result)
	if getReqErr != nil {
		return nil, getReqErr
	}

	return &result, nil
}

func (c *client) UnblockUser(userID, unblockUserID int) error {
	pathString := fmt.Sprintf("/users/%d/block/%d", userID, unblockUserID)
	parsedURL := c.PrepareUrl(pathString)

	result := struct{}{}

	getReqErr := c.delete(parsedURL, "", &result)
	if getReqErr != nil {
		return getReqErr
	}

	return nil
}
