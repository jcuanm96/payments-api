package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/auth"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/kevinburke/twilio-go"
)

func CheckPhone(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	req := incomeRequest.(request.Phone)

	res, err := svc.CheckPhone(ctx, req)

	return res, err
}

func CheckUsername(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	req := incomeRequest.(request.Username)

	res, err := svc.CheckUsername(ctx, req)

	return res, err
}

func CheckEmail(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	req := incomeRequest.(request.Email)

	res, err := svc.CheckEmail(ctx, req)

	return res, err
}

func CheckGoatInviteCode(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	req := incomeRequest.(request.GoatInviteCode)

	res, err := svc.CheckGoatInviteCode(ctx, req)

	return res, err
}

func VerifySMS(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	req := incomeRequest.(request.VerifySMS)

	res, err := svc.VerifySMS(ctx, req)

	return res, err
}

func SignOut(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	req := incomeRequest.(request.SignOut)

	err := svc.SignOut(ctx, req)

	return nil, err
}

func SignInSMS(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	req := incomeRequest.(request.SignInSMS)

	res, err := svc.SignInSMS(ctx, req)

	return res, err
}

func SignupSMS(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	req := incomeRequest.(request.SignupSMS)

	res, err := svc.SignupSMS(ctx, req)

	return res, err
}

func SignupSMSGoat(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	req := incomeRequest.(request.SignupSMSGoat)

	formattedNumber, formatErr := req.FormattedPhonenumber()
	if formatErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to sign up. Please try again.",
			fmt.Sprintf("Could not format phone number %s: %v", req.Number, formatErr),
		)
	}

	// Ignore Twilio code check if we're testing
	phonenumberCheckStatus := twilio.CheckPhoneNumber{
		Status: "approved",
	}
	if appconfig.Config.Gcloud.Project == "vama-prod" {
		status, twilioCheckErr := svc.VerifyTwilioCode(ctx, formattedNumber, req.Code)
		if twilioCheckErr != nil {
			vlog.Errorf(ctx, "Invalid verification code. Err: %v", twilioCheckErr)
			return nil, httperr.NewCtx(
				ctx,
				400,
				http.StatusBadRequest,
				"Invalid verification code.",
			)
		}
		phonenumberCheckStatus = *status
	}

	if phonenumberCheckStatus.Status != "approved" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Invalid verification code.",
			"Invalid verification code.",
		)
	}

	res, err := svc.SignupSMSGoat(ctx, formattedNumber, req)

	return res, err
}

func RefreshToken(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	req := incomeRequest.(request.RefreshToken)

	res, err := svc.RefreshToken(ctx, req)

	return res, err
}

func GenerateGoatInviteCode(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(auth.Usecase)

	shouldInsertCode := true
	res, err := svc.GenerateGoatInviteCode(ctx, shouldInsertCode)

	return res, err
}
