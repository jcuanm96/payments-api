package test_util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func MakeGetRequest(t *testing.T, app *fiber.App, endpoint string, statusCode int, params map[string]string, asUserID int, resp interface{}) {
	ctx := context.Background()

	req, httpErr := http.NewRequest("GET", endpoint, nil)
	if httpErr != nil {
		vlog.Errorf(ctx, "Error constructing http request: %s", httpErr.Error())
		t.Fail()
	}
	req.Header.Set("id", fmt.Sprint(asUserID))
	q := req.URL.Query()
	for key := range params {
		q.Add(key, params[key])
	}
	req.URL.RawQuery = q.Encode()

	httpResp, reqErr := app.Test(req, -1)
	require.Nil(t, reqErr)
	require.NotNil(t, httpResp)
	require.Equal(t, statusCode, httpResp.StatusCode)

	defer httpResp.Body.Close()
	jsonErr := json.NewDecoder(httpResp.Body).Decode(resp)
	require.Nil(t, jsonErr)
}

func MakeGetRequestAssert200(t *testing.T, app *fiber.App, endpoint string, params map[string]string, asUserID int, resp interface{}) {
	MakeGetRequest(t, app, endpoint, 200, params, asUserID, resp)
}

func MakePostRequest(t *testing.T, app *fiber.App, endpoint string, statusCode int, params map[string]interface{}, asUserID *int, resp interface{}) {
	ctx := context.Background()

	requestBody, reqBodyErr := json.Marshal(params)
	if reqBodyErr != nil {
		vlog.Errorf(ctx, "Error marshaling post request body params: %s", reqBodyErr.Error())
		t.Fail()
	}

	req, httpErr := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if httpErr != nil {
		vlog.Errorf(ctx, "Error constructing http request: %s", httpErr.Error())
		t.Fail()
	}

	if asUserID != nil {
		req.Header.Set("id", fmt.Sprint(*asUserID))
	}

	req.Header.Set("Content-Type", "application/json")
	httpResp, reqErr := app.Test(req, -1)
	require.Nil(t, reqErr)
	require.NotNil(t, httpResp)
	require.Equal(t, statusCode, httpResp.StatusCode)

	defer httpResp.Body.Close()
	jsonErr := json.NewDecoder(httpResp.Body).Decode(resp)
	require.Nil(t, jsonErr)
}

func MakePostRequestAssert200(t *testing.T, app *fiber.App, endpoint string, params map[string]interface{}, asUserID *int, resp interface{}) {
	MakePostRequest(t, app, endpoint, 200, params, asUserID, resp)
}

func MakePostRequestAssert403(t *testing.T, app *fiber.App, endpoint string, params map[string]interface{}, asUserID *int, resp interface{}) {
	MakePostRequest(t, app, endpoint, 403, params, asUserID, resp)
}

func MakePatchRequest(t *testing.T, app *fiber.App, endpoint string, statusCode int, params map[string]interface{}, asUserID *int, resp interface{}) {
	ctx := context.Background()

	requestBody, reqBodyErr := json.Marshal(params)
	if reqBodyErr != nil {
		vlog.Errorf(ctx, "Error marshaling patch request body params: %s", reqBodyErr.Error())
		t.Fail()
	}

	req, httpErr := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(requestBody))
	if httpErr != nil {
		vlog.Errorf(ctx, "Error constructing http request: %s", httpErr.Error())
		t.Fail()
	}

	if asUserID != nil {
		req.Header.Set("id", fmt.Sprint(*asUserID))
	}

	req.Header.Set("Content-Type", "application/json")
	httpResp, reqErr := app.Test(req, -1)
	require.Nil(t, reqErr)
	require.NotNil(t, httpResp)
	require.Equal(t, statusCode, httpResp.StatusCode)

	defer httpResp.Body.Close()
	jsonErr := json.NewDecoder(httpResp.Body).Decode(resp)
	require.Nil(t, jsonErr)
}

func MakePatchRequestAssert200(t *testing.T, app *fiber.App, endpoint string, params map[string]interface{}, asUserID *int, resp interface{}) {
	MakePatchRequest(t, app, endpoint, 200, params, asUserID, resp)
}

func MakeDeleteRequest(t *testing.T, app *fiber.App, endpoint string, statusCode int, params map[string]interface{}, asUserID int, resp interface{}) {
	ctx := context.Background()

	requestBody, reqBodyErr := json.Marshal(params)
	if reqBodyErr != nil {
		vlog.Errorf(ctx, "Error marshaling delete request body params: %s", reqBodyErr.Error())
		t.Fail()
	}

	req, httpErr := http.NewRequest("DELETE", endpoint, bytes.NewBuffer(requestBody))
	if httpErr != nil {
		vlog.Errorf(ctx, "Error constructing http request: %s", httpErr.Error())
		t.Fail()
	}
	req.Header.Set("id", fmt.Sprint(asUserID))
	req.Header.Set("Content-Type", "application/json")

	httpResp, reqErr := app.Test(req, -1)
	require.Nil(t, reqErr)
	require.NotNil(t, httpResp)
	require.Equal(t, statusCode, httpResp.StatusCode)

	defer httpResp.Body.Close()
	jsonErr := json.NewDecoder(httpResp.Body).Decode(resp)
	require.Nil(t, jsonErr)
}

func MakeDeleteRequestAssert200(t *testing.T, app *fiber.App, endpoint string, params map[string]interface{}, asUserID int, resp interface{}) {
	MakeDeleteRequest(t, app, endpoint, 200, params, asUserID, resp)
}
