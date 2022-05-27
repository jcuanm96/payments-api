package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const listUnpaidProvidersDefaultErr = "Something went wrong listing unpaid providers."

func (svc *usecase) ListUnpaidProviders(ctx context.Context) (*response.ListUnpaidProviders, error) {
	providers, listProvidersErr := svc.repo.ListUnpaidProviders(ctx)
	if listProvidersErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			listUnpaidProvidersDefaultErr,
			fmt.Sprintf("Error listing unpaid providers: %v", listProvidersErr),
		)
	}
	return &response.ListUnpaidProviders{Providers: providers}, nil
}
