package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) SearchUsersExecQuery(ctx context.Context, query string, params []string) ([]response.User, error) {
	res := make([]response.User, 0)

	newParams := make([]interface{}, len(params))
	for i, v := range params {
		newParams[i] = v
	}
	rows, queryErr := s.MasterNode().Query(ctx, query, newParams...)

	if queryErr != nil {
		if queryErr == pgx.ErrNoRows {
			return res, nil
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, newParams))
		return res, queryErr
	}

	defer rows.Close()

	for rows.Next() {
		item := response.User{}
		scanErr := rows.Scan(
			&item.ID,
			&item.UUID,
			&item.FirstName,
			&item.LastName,
			&item.Phonenumber,
			&item.CountryCode,
			&item.Email,
			&item.Type,
			&item.ProfileAvatar,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Username,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		res = append(res, item)
	}
	return res, nil
}
