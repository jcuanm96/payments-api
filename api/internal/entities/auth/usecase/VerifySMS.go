package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) VerifySMS(ctx context.Context, req request.VerifySMS) (*response.AuthConfirm, error) {
	formattedNumber, formatErr := req.FormattedPhonenumber()
	if formatErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong.  Please try again.",
			fmt.Sprintf("Could not format phone number %s. Err: %v", req.Number, formatErr.Error()),
		)
	}
	res := response.AuthConfirm{}
	_, twilioErr := svc.twilio.Create(ctx, appconfig.Config.Twilio.Verify, url.Values{"To": []string{formattedNumber}, "Channel": []string{"sms"}})
	if twilioErr != nil {
		if strings.Contains(twilioErr.Error(), "landline") {
			return nil, httperr.NewCtx(
				ctx,
				400,
				http.StatusBadRequest,
				"Cannot send verification code to landlines.",
				"Cannot send verification code to landlines.",
			)
		} else if strings.Contains(twilioErr.Error(), "Max") && strings.Contains(twilioErr.Error(), "attempt") {
			return nil, httperr.NewCtx(
				ctx,
				400,
				http.StatusBadRequest,
				"Max send attempts reached. Please try again in 10 minutes.",
				"Max send attempts reached. Please try again in 10 minutes.",
			)
		} else {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
				fmt.Sprintf("Error trying to get verification code for phone number %s. Err: %v", formattedNumber, twilioErr),
			)
		}
	}
	return &res, nil
}
