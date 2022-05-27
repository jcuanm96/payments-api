package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetUsers(ctx context.Context, query string, userID int, limit int) ([]response.User, error) {
	gridListReq := request.NewGridList{
		PageSize:   limit,
		PageNumber: 1,
	}
	gridListReq.CustomFilters = make([]request.CustomFilterItem, 0)
	gridListReq.CustomFilters = append(gridListReq.CustomFilters, request.CustomFilterItem{
		Name:   "userTypes",
		Values: []string{"USER", "GOAT"},
	})
	gridListReq.CustomFilters = append(gridListReq.CustomFilters, request.CustomFilterItem{
		Name:   "query",
		Values: []string{query},
	})
	gridListReq.Sorts = make([]request.SortItem, 0)
	gridListReq.Sorts = append(gridListReq.Sorts, request.SortItem{
		Field: "first_name",
		Dir:   "asc",
	})
	gridListReq.Sorts = append(gridListReq.Sorts, request.SortItem{
		Field: "last_name",
		Dir:   "asc",
	})
	gridListReq.Sorts = append(gridListReq.Sorts, request.SortItem{
		Field: "id",
		Dir:   "asc",
	})
	_, users, searchUsersErr := svc.SearchUsers(ctx, userID, gridListReq)

	if searchUsersErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to search for users.  Please try again.",
			fmt.Sprintf("Error retrieving query results for user %d: %v", userID, searchUsersErr),
		)
	}

	if users == nil || len(users) < 1 {
		return []response.User{}, nil
	}

	return users, nil
}
