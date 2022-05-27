package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	userrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) CheckLinkSuffixIsTaken(ctx context.Context, req request.CheckLinkSuffixIsTaken) (*response.CheckLinkSuffixIsTaken, error) {
	taken, isTakenErr := userrepo.IsLinkSuffixTaken(ctx, svc.repo.MasterNode(), req.LinkSuffix)
	if isTakenErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to check link.",
			fmt.Sprintf("Error checking if link is taken: %v", isTakenErr),
		)
	}

	isTaken := taken.IsTaken

	if req.ChannelID != nil {
		// if link is taken but it's taken by this channel, report as not taken.
		if isTaken &&
			((taken.PaidGroupChannelID != nil && *req.ChannelID == *taken.PaidGroupChannelID) ||
				(taken.FreeGroupChannelID != nil && *req.ChannelID == *taken.FreeGroupChannelID)) {
			isTaken = false
		}
	}

	return &response.CheckLinkSuffixIsTaken{Taken: isTaken}, nil
}
