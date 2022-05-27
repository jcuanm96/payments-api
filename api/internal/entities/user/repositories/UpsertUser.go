package repositories

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func UpsertUserDB(ctx context.Context, tx pgx.Tx, item *response.User) (*response.User, error) {
	query, args, squirrelErr := squirrel.Insert("core.users").
		Columns(
			"first_name",
			"last_name",
			"phone_number",
			"country_code",
			"email",
			"username",
			"user_type",
			"profile_avatar",
			"stripe_id",
		).
		Values(
			item.FirstName,
			item.LastName,
			item.Phonenumber,
			item.CountryCode,
			item.Email,
			item.Username,
			item.Type,
			item.ProfileAvatar,
			item.StripeID,
		).
		Suffix(`
		ON conflict (phone_number)
			DO update SET 
			first_name = ?,
			last_name = ?,
			phone_number = ?,
			country_code = ?,
			email = ?,
			username = ?,
			user_type = ?,
			profile_avatar = ?,
			stripe_id = ?`,
			item.FirstName,
			item.LastName,
			item.Phonenumber,
			item.CountryCode,
			item.Email,
			item.Username,
			item.Type,
			item.ProfileAvatar,
			item.StripeID).
		Suffix(`
		RETURNING
			id,
			uuid,
			first_name,
			last_name,
			phone_number,
			country_code,
			email,
			username,
			user_type,
			profile_avatar,
			stripe_id
		`).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(squirrelErr, query, args))
		return nil, squirrelErr
	}

	row, sqlErr := tx.Query(ctx, query, args...)
	if sqlErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(sqlErr, query, args))
		return nil, sqlErr
	}

	defer row.Close()
	if !row.Next() {
		return nil, fmt.Errorf("no rows returned for upserted user")
	}

	res := response.User{}
	scanErr := row.Scan(
		&res.ID,
		&res.UUID,
		&res.FirstName,
		&res.LastName,
		&res.Phonenumber,
		&res.CountryCode,
		&res.Email,
		&res.Username,
		&res.Type,
		&res.ProfileAvatar,
		&res.StripeID,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	return &res, nil
}

func (s *repository) UpsertUser(ctx context.Context, tx pgx.Tx, item *response.User) (*response.User, error) {
	return UpsertUserDB(ctx, tx, item)
}
