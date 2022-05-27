package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func CalcPages(s BaseRepository, ctx context.Context, query string, params []string, reqPageNumber, reqPageSize, currentCount int) (response.Paging, error) {
	res := response.Paging{}
	res.CurrentCount = currentCount
	res.Page = reqPageNumber
	res.Size = reqPageSize
	return res, nil
}
