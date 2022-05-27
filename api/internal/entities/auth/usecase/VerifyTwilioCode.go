package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	twilio "github.com/kevinburke/twilio-go"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) VerifyTwilioCode(ctx context.Context, phone string, code string) (*twilio.CheckPhoneNumber, error) {
	status, twilioCheckErr := svc.twilio.Check(ctx, appconfig.Config.Twilio.Verify, url.Values{"To": []string{phone}, "Code": []string{code}})
	if twilioCheckErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Invalid verification code.",
			fmt.Sprintf("Invalid verification code for phone %s and code %s. Err: %v.", phone, code, twilioCheckErr),
		)
	} else if status == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when validating code.",
			fmt.Sprintf("status came back as nil for phone %s and code %s.", phone, code),
		)
	}

	return status, nil
}
