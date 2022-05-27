package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetBioData(ctx context.Context, req request.GetBioData) (*response.GetBioData, error) {
	runnable := svc.repo.MasterNode()
	user, checkUsernameErr := svc.user.CheckUsernameAlreadyExists(ctx, runnable, req.Username)
	if checkUsernameErr != nil {
		vlog.Errorf(ctx, "Could not find user for name: %s: %v", req.Username, checkUsernameErr)
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
		)
	} else if user == nil {
		vlog.Errorf(ctx, "Could not get bio data for user %s because they do not exist", req.Username)
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			constants.ErrSomethingWentWrong,
		)
	}

	res, getBioDataErr := svc.repo.GetBioData(ctx, runnable, strings.ToLower(req.Username), user.ID)
	if getBioDataErr != nil {
		vlog.Errorf(ctx, "Could not get bio data for user %s: %v", req.Username, getBioDataErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	}

	return res, nil
}
