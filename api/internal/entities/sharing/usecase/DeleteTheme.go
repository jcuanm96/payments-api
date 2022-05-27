package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) DeleteTheme(ctx context.Context, req request.DeleteTheme) error {
	deleteThemeErr := svc.repo.DeleteTheme(ctx, req.ThemeID)
	if deleteThemeErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Could not delete theme with ID %d: %v", req.ThemeID, deleteThemeErr),
		)
	}

	return nil
}
