package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (s *repository) IsContact(ctx context.Context, userID, contactID int) (bool, error) {
	return IsContact(ctx, s.MasterNode(), userID, contactID)
}

func IsContact(ctx context.Context, runnable utils.Runnable, userID, contactID int) (bool, error) {
	query, args, squirrelErr := squirrel.Select(
		"id",
	).
		From("core.users_contacts").
		Where("user_id = ?", userID).
		Where("contact_id = ?", contactID).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return false, squirrelErr
	}

	rows, queryErr := runnable.Query(ctx, query, args...)
	if queryErr != nil {
		return false, queryErr
	}

	defer rows.Close()
	return rows.Next(), nil
}
