package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func UpdateUser(ctx context.Context, tx pgx.Tx, phonenumber string, email string) error {
	query, args, squirrelErr := squirrel.Update("core.users").
		Set("user_type", "GOAT").
		Set("email", email).
		Where("phone_number=?", phonenumber).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := tx.Exec(ctx, query, args...)
	if queryErr != nil {
		return queryErr
	}
	return nil
}
