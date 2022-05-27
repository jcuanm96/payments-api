package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func UpsertVamaUser(tx pgx.Tx, vamaUser *response.User) error {
	query, args, squirrelErr := squirrel.Insert("core.users").
		Columns(
			"id",
			"first_name",
			"last_name",
			"phone_number",
			"email",
			"username",
			"user_type",
			"stripe_id",
		).
		Values(
			vamaUser.ID,
			vamaUser.FirstName,
			vamaUser.LastName,
			vamaUser.Phonenumber,
			vamaUser.Email,
			vamaUser.Username,
			vamaUser.Type,
			vamaUser.StripeID,
		).
		Suffix(`
	ON conflict (id)
		DO update SET 
		first_name = ?,
		last_name = ?,
		phone_number = ?,
		email = ?,
		username = ?,
		user_type = ?,
		stripe_id = ?`,
			vamaUser.FirstName,
			vamaUser.LastName,
			vamaUser.Phonenumber,
			vamaUser.Email,
			vamaUser.Username,
			vamaUser.Type,
			vamaUser.StripeID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := tx.Exec(context.Background(), query, args...)
	if queryErr != nil {
		logrus.Error(utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}
	return nil
}
