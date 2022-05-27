package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errSearchingUsers = "Something went wrong when trying to search for users.  Please try again."

func (svc *usecase) ListUsers(ctx context.Context, req request.ListUsers) (*response.ListUsers, error) {
	user, getCurrentUserErr := svc.GetCurrentUser(ctx)
	if getCurrentUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errSearchingUsers,
			fmt.Sprintf("Error getting current user: %v", getCurrentUserErr),
		)
	} else if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errSearchingUsers,
			"user was nil in ListUsers",
		)
	}

	gridListReq := request.NewGridList{
		PageSize:   req.Size,
		PageNumber: req.Page,
	}
	gridListReq.CustomFilters = make([]request.CustomFilterItem, 0)
	gridListReq.CustomFilters = append(gridListReq.CustomFilters, request.CustomFilterItem{
		Name:   "userTypes",
		Values: getUserTypes(req.UserTypes),
	})
	{
		customFilters := strings.Split(req.Filters, ",")
		for _, v := range customFilters {
			if v != "" {
				gridListReq.CustomFilters = append(gridListReq.CustomFilters, request.CustomFilterItem{
					Name:   v,
					Values: []string{"true"},
				})
			}
		}
	}
	gridListReq.CustomFilters = append(gridListReq.CustomFilters, request.CustomFilterItem{
		Name:   "query",
		Values: []string{req.Query},
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
	paging, users, searchUsersErr := (*svc.search).SearchUsers(ctx, user.ID, gridListReq)

	if searchUsersErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errSearchingUsers,
			fmt.Sprintf("Error retrieving query results: %v", searchUsersErr),
		)
	}

	if users == nil {
		return &response.ListUsers{}, nil
	}

	res := response.ListUsers{
		Paging: *paging,
		Users:  users,
	}
	return &res, nil
}

func getUserTypes(t string) []string {
	if t == "" {
		return []string{"GOAT", "USER"}
	}
	types := strings.Split(t, ",")
	res := make([]string, 0)
	for _, typ := range types {
		typ = strings.ToUpper(strings.TrimSpace(typ))
		res = append(res, typ)
	}
	return res
}
