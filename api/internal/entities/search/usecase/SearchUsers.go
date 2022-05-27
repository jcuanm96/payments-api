package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (svc *usecase) SearchUsers(ctx context.Context, userID int, req request.NewGridList) (*response.Paging, []response.User, error) {
	query, sortDef, params, err := svc.repo.SearchUsersCalcQuery(ctx, userID, req.CustomFilters)
	if err != nil {
		return nil, nil, err
	}
	queryWithSortAndOffset, err := svc.repo.CalcQueryWithSortAndOffset(ctx, query, req.Sorts, req.PageSize, req.PageNumber, sortDef)
	if err != nil {
		return nil, nil, err
	}

	resItems, err := svc.repo.SearchUsersExecQuery(ctx, queryWithSortAndOffset, params)
	if err != nil {
		return nil, nil, err
	}
	paging, err := svc.repo.CalcPages(ctx, query, params, req.PageNumber, req.PageSize, len(resItems))
	if err != nil {
		return nil, nil, err
	}
	return &paging, resItems, nil
}
