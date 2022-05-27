package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getGoatUsersErr = "Something went wrong when searching for creators. Please try again."

func (svc *usecase) GetGoatUsers(ctx context.Context, req request.GetGoatUsers) (*response.GetGoatUsers, error) {
	var currUserID int
	if req.ExcludeSelf {
		currUser, userErr := svc.GetCurrentUser(ctx)
		if userErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				getGoatUsersErr,
				fmt.Sprintf("Error finding user in the current context: %v", userErr),
			)

		}
		if currUser == nil {
			return nil, httperr.NewCtx(
				ctx,
				404,
				http.StatusNotFound,
				getGoatUsersErr,
				"currUser was nil in GetGoatUsers.",
			)
		}
		currUserID = currUser.ID
	}

	goatUsers, getGoatUsersRepoErr := svc.repo.GetGoatUsers(ctx, currUserID, req)
	if getGoatUsersRepoErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getGoatUsersErr,
			fmt.Sprintf("Error getting creator users: %v", getGoatUsersRepoErr),
		)
	}

	res := response.GetGoatUsers{
		GoatUsers: goatUsers,
	}

	return &res, nil
}
