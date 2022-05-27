package repositories

import (
	"context"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func GetCountByQuery(s BaseRepository, ctx context.Context, query string, params []string) (int, error) {
	res := 0
	newParams := make([]interface{}, len(params))
	for i, v := range params {
		newParams[i] = v
	}
	q := fmt.Sprintf("select coalesce(count(*),0) from (%s) t", query)
	rows, queryErr := s.MasterNode().Query(ctx, q, newParams...)

	if queryErr != nil {
		if queryErr == pgx.ErrNoRows {
			return res, nil
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, q, nil))
		return res, queryErr
	}

	defer rows.Close()

	for rows.Next() {
		item := 0
		scanErr := rows.Scan(
			&item,
		)
		if scanErr != nil {
			return 0, scanErr
		}
		res = item
	}
	return res, nil
}
