package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/nyaruka/phonenumbers"
)

func (svc *usecase) CheckPhone(ctx context.Context, req request.Phone) (*response.Check, error) {
	res := response.Check{}

	if !validatePhonenumber(req.Number, req.CountryCode) {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Please enter a valid phone number.",
			"Please enter a valid phone number.",
		)
	}
	formattedNumber, formatErr := req.FormattedPhonenumber()
	if formatErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when checking that the provided phone number is available.",
			fmt.Sprintf("Could not format phone number %s. Err: %v", req.Number, formatErr),
		)
	}
	taken, err := svc.user.CheckPhoneAlreadyExists(ctx, formattedNumber)

	if err != nil {
		return nil, err
	}
	if taken {
		res.IsTaken = true
		res.Message = "A user with this phone number already exists."
	}
	return &res, nil
}

func validatePhonenumber(val, countryCode string) bool {
	num, err := phonenumbers.Parse(val, countryCode)
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(num)
}
