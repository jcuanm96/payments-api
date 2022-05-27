package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getGoatChatPriceErr = "Something went wrong when getting creator chat price."

func (svc *usecase) GetGoatChatPrice(ctx context.Context, req request.GetGoatChatPrice) (*response.GetGoatChatPrice, error) {
	res, repoErr := svc.repo.GetGoatChatPrice(ctx, req.GoatUserID)
	if repoErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getGoatChatPriceErr,
			fmt.Sprintf("Something went wrong when getting creator chat price. Err: %v", repoErr),
		)
	}
	if res == nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			getGoatChatPriceErr,
			"Response was nil when getting creator chat price.",
		)
	}
	return res, nil
}
