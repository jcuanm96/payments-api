package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errReportingUser = "Something went wrong when trying to report user. Please try again."

func (svc *usecase) ReportUser(ctx context.Context, req request.ReportUser) error {
	user, userErr := svc.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errReportingUser,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	} else if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errReportingUser,
			"user was nil in ReportUser",
		)
	}

	reportUserErr := svc.repo.ReportUser(ctx, user.ID, req.ReportedUserID, req.Description)
	if reportUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errReportingUser,
			fmt.Sprintf("Error inserting report entry into db: %v", reportUserErr),
		)
	}

	return nil
}
