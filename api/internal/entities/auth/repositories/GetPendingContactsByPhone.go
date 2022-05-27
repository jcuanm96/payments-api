package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetPendingContactsByPhone(ctx context.Context, phone string) ([]int, error) {
	query, args, squirrelErr := squirrel.Select(
		"user_id",
	).
		From("core.pending_contacts").
		Where("phone_number = ?", phone).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}

	defer rows.Close()
	alreadyExistingContactsIDs := []int{}

	for rows.Next() {
		var userID int
		scanErr := rows.Scan(
			&userID,
		)
		if scanErr != nil {
			return nil, scanErr
		}

		alreadyExistingContactsIDs = append(alreadyExistingContactsIDs, userID)
	}

	return alreadyExistingContactsIDs, nil
}
