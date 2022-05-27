package service

import (
	"context"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetMinRequiredVersionByDevice(ctx context.Context, req request.GetMinRequiredVersionByDevice) (*response.GetMinRequiredVersionByDevice, error) {
	validPlatforms := map[string]struct{}{
		"ios": {},
	}

	if _, ok := validPlatforms[req.Platform]; !ok {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Invalid platform was passed.",
		)
	}

	res := response.GetMinRequiredVersionByDevice{
		Version: constants.MIN_REQUIRED_IOS_APP_VERSION,
	}
	return &res, nil
}
