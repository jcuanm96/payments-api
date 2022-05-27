package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const listPayoutHistoryDefaultErr = "Something went wrong listing payout history."

func (svc *usecase) ListPayoutHistory(ctx context.Context, req request.ListPayoutHistory) (*response.ListPayoutHistory, error) {
	limit := req.Limit + 1
	if req.CursorPayPeriodID == 0 {
		currPayPeriod, getPayPeriodErr := svc.repo.GetPayPeriodPerTS(ctx, time.Now().Unix())
		if getPayPeriodErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				listPayoutHistoryDefaultErr,
				fmt.Sprintf("Error getting pay period in ListPayoutHistory: %v", getPayPeriodErr),
			)
		} else if currPayPeriod == nil {
			return nil, httperr.NewCtx(
				ctx,
				404,
				http.StatusNotFound,
				"The current pay period does not exist.",
				"The current pay period does not exist.",
			)
		}

		req.CursorPayPeriodID = int64(currPayPeriod.ID) // id of current pay period, which is non-inclusive
	}

	history, listPayoutHistoryErr := svc.repo.ListPayoutHistory(ctx, req, limit)
	if listPayoutHistoryErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			listPayoutHistoryDefaultErr,
			fmt.Sprintf("Error listing payout history: %v", listPayoutHistoryErr),
		)
	}

	hasNext := len(history) > int(req.Limit)

	if len(history) > int(req.Limit) {
		history = history[:req.Limit]
	}

	res := response.ListPayoutHistory{
		PayoutHistory: history,
		HasNext:       hasNext,
	}

	return &res, nil
}
