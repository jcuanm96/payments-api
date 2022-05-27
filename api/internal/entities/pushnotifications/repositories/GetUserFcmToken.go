package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
)

func (s *repository) GetUserFcmToken(ctx context.Context, userID int) (string, error) {
	query, args, squirrelErr := squirrel.Select(
		"COALESCE(fcm_token, '')",
	).
		From("core.tokens").
		Where("user_id = ?", userID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return "", squirrelErr
	}

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		return "", queryErr
	}

	defer rows.Close()

	if !rows.Next() {
		return "", constants.ErrNotFound
	}

	var token string
	scanErr := rows.Scan(
		&token,
	)
	if scanErr != nil {
		return "", scanErr
	}

	return token, nil
}
