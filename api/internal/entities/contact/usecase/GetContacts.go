package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getContactsErr = "Something went wrong when trying to get contacts."

func (svc *usecase) GetContacts(ctx context.Context, req request.GetContacts) (*response.GetContacts, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getContactsErr,
			fmt.Sprintf("Error getting current user: %s", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusInternalServerError,
			getContactsErr,
			"user came back nil in GetContacts",
		)
	}

	gridList := request.NewGridList{
		PageSize:   req.Size,
		PageNumber: req.Page,
	}

	gridList.CustomFilters = make([]request.CustomFilterItem, 0)
	gridList.CustomFilters = append(gridList.CustomFilters, request.CustomFilterItem{
		Name:   "onlyContacts",
		Values: []string{"true"},
	})
	gridList.CustomFilters = append(gridList.CustomFilters, request.CustomFilterItem{
		Name:   "query",
		Values: []string{""},
	})

	if req.ExcludeGoats {
		gridList.CustomFilters = append(gridList.CustomFilters, request.CustomFilterItem{
			Name:   "userTypes",
			Values: []string{"USER"},
		})
	}

	gridList.Sorts = make([]request.SortItem, 0)
	gridList.Sorts = append(gridList.Sorts, request.SortItem{
		Field: "first_name",
		Dir:   "asc",
	})
	gridList.Sorts = append(gridList.Sorts, request.SortItem{
		Field: "last_name",
		Dir:   "asc",
	})
	gridList.Sorts = append(gridList.Sorts, request.SortItem{
		Field: "id",
		Dir:   "asc",
	})

	paging, contacts, searchErr := svc.search.SearchUsers(ctx, user.ID, gridList)

	res := response.GetContacts{
		Paging:   paging,
		Contacts: contacts,
	}

	if searchErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getContactsErr,
			fmt.Sprintf("Error searching for contact.  Err: %s", searchErr),
		)
	}
	return &res, nil
}
