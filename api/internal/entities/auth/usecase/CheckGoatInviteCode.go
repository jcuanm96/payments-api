package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) CheckGoatInviteCode(ctx context.Context, req request.GoatInviteCode) (*response.Check, error) {
	res := response.Check{}

	userID, wasFound, getStatusErr := svc.user.GetInviteCodeStatus(ctx, req.Code)
	if getStatusErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when checking your invite code.  Please try again.",
			fmt.Sprintf("Error checking Goat invite code status for code %s. Error: %v", req.Code, getStatusErr),
		)
	}

	if wasFound && userID == nil {
		return &res, nil
	} else {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"The invite code you entered is either invalid or taken.",
			"The invite code you entered is either invalid or taken.",
		)
	}
}
