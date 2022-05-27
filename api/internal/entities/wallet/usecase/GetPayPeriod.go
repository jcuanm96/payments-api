package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getPayoutPeriodDefaultErr = "Error when getting payout period"

func (svc *usecase) GetPayPeriod(ctx context.Context, req request.GetPayPeriod) (*response.PayoutPeriod, error) {
	payPeriod, getPayPeriodErr := svc.repo.GetPayPeriodPerTS(ctx, req.Timestamp)
	if getPayPeriodErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getPayoutPeriodDefaultErr,
			fmt.Sprintf("Error getting pay period: %v", getPayPeriodErr),
		)
	} else if payPeriod == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"The pay period that was inputted does not exist.",
			"The pay period that was inputted does not exist.",
		)
	}

	return payPeriod, nil
}
