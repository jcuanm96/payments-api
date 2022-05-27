package controller

import (
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
)

func ComposeEndpoints() []sconfiguration.EndpointConfiguration {

	res := []sconfiguration.EndpointConfiguration{
		{
			Path:           "/sign-up/sms",
			Method:         "post",
			Handler:        SignupSMS,
			RequestDecoder: decodeSignupSMS,
		},
		{
			Path:           "/sign-up/sms/goat",
			Method:         "post",
			Handler:        SignupSMSGoat,
			RequestDecoder: decodeSignupSMSGoat,
		},
		{
			Path:           "/verify/sms/resend",
			Method:         "get",
			Handler:        VerifySMS,
			RequestDecoder: decodeVerifySMS,
		},
		{
			Path:           "/sign-in/sms",
			Method:         "post",
			Handler:        SignInSMS,
			RequestDecoder: decodeSignInSMS,
		},
		{
			Path:           "/verify/sms",
			Method:         "get",
			Handler:        VerifySMS,
			RequestDecoder: decodeVerifySMS,
		},
		{
			Path:           "/refresh",
			Method:         "post",
			Handler:        RefreshToken,
			RequestDecoder: decodeRefreshToken,
		},
		{
			Path:            "/sign-out",
			Method:          "post",
			Handler:         SignOut,
			RequestDecoder:  decodeSignOut,
			ResponseEncoder: sconfiguration.SuccessResponse,
		},
		{
			Path:           "/check/email",
			Method:         "get",
			Handler:        CheckEmail,
			RequestDecoder: decodeCheckEmail,
		},
		{
			Path:           "/check/invite-code/goat",
			Method:         "get",
			Handler:        CheckGoatInviteCode,
			RequestDecoder: decodeCheckGoatInviteCode,
		},
		{
			Path:           "/check/phone",
			Method:         "get",
			Handler:        CheckPhone,
			RequestDecoder: decodeCheckPhone,
		},
		{
			Path:           "/check/username",
			Method:         "get",
			Handler:        CheckUsername,
			RequestDecoder: decodeCheckUsername,
		},
		{
			Path:           "/invite-code/goat",
			Method:         "post",
			Handler:        GenerateGoatInviteCode,
			RequestDecoder: decodeGenerateGoatInviteCode,
		},
	}

	return res
}
