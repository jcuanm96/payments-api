package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) ListBannedUsers(ctx context.Context, req request.ListBannedUsers) (*response.ListBannedUsers, error) {
	res, listBannedUsersErr := svc.repo.ListBannedUsers(ctx, svc.repo.MasterNode(), req)
	if listBannedUsersErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong getting banned users.",
			fmt.Sprintf("Error listing banned users: %v", listBannedUsersErr),
		)
	}
	return res, nil
}
