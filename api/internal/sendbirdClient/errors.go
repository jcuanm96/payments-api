package sendbird

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Codes
const ResourceNotFound = 400201

// Messages
var ErrGroupChannelNotFound = errors.New("sendbird group channel does not exist")

type SendbirdErrorResponse struct {
	HasError bool   `json:"error"`
	Message  string `json:"message"`
	Code     int    `json:"code"`
}

// Implement error interface
func (s SendbirdErrorResponse) Error() string {
	if s.Code != 200 && s.Code != 0 {
		return fmt.Sprintf("SendbirdError: %d - %s", s.Code, s.Message) // or s.message or some kind of format
	}
	return "{}"
}

func CheckSendbirdError(httpResp *http.Response) error {
	if httpResp.StatusCode != 200 {
		errorMessageBody := SendbirdErrorResponse{}
		err := json.NewDecoder(httpResp.Body).Decode(&errorMessageBody)
		if err != nil {
			return fmt.Errorf("sendbird client error: %s", err)
		}

		return errorMessageBody
	}
	return nil
}
