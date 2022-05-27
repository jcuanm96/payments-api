package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getPayoutPeriodsDefaultErr = "Error when getting payout periods"

func (svc *usecase) GetPayoutPeriods(ctx context.Context, req request.GetPayoutPeriods) (*response.GetPayoutPeriods, error) {
	payoutPeriods, getPayoutPeriodsErr := svc.repo.GetPayoutPeriods(ctx, req)
	if getPayoutPeriodsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getPayoutPeriodsDefaultErr,
			fmt.Sprintf("Error when getting payout periods. Err: %v", getPayoutPeriodsErr),
		)
	}

	res := response.GetPayoutPeriods{
		PayoutPeriods: payoutPeriods,
	}

	return &res, nil
}
