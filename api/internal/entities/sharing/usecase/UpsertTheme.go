package service

import (
	"context"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) UpsertTheme(ctx context.Context, req request.UpsertTheme) (*response.UpsertTheme, error) {
	upsertErr := svc.repo.UpsertTheme(ctx, req)

	if upsertErr != nil {
		vlog.Errorf(ctx, "Error occurred when upserting theme %s to db: %v", req.Name, upsertErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong upserting theme. Please try again.",
		)
	}

	res := response.UpsertTheme{}

	return &res, nil
}
