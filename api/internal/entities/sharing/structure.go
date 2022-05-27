package sharing

import (
	"context"
	"net/url"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
)

// https://firebase.google.com/docs/dynamic-links/create-manually
func CreateDynamicLink(ctx context.Context, deepLink string) (string, error) {
	dynamicLinkParams := map[string]string{
		"link": deepLink,
		"isi":  constants.APPLE_APP_STORE_ID,
		"ibi":  appconfig.Config.Gcloud.AppleBundleID,
		"efr":  "1",
	}

	dynamicLinkBase, urlParseErr := url.Parse(appconfig.Config.Gcloud.DynamicLinkBaseURL)
	if urlParseErr != nil {
		return "", urlParseErr
	}

	params := url.Values{}
	for key, value := range dynamicLinkParams {
		params.Add(key, value)
	}

	dynamicLinkBase.RawQuery = params.Encode()
	return dynamicLinkBase.String(), nil
}
