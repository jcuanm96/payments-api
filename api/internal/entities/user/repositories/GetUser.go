package repositories

import (
	"context"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetUser(ctx context.Context, runnable utils.Runnable, prefict string, params ...interface{}) (*response.User, error) {
	res := response.User{}

	query, args, squirrelErr := squirrel.Select(
		"id",
		"uuid",
		"stripe_id",
		"first_name",
		"last_name",
		"phone_number",
		"country_code",
		"email",
		"username",
		"user_type",
		"profile_avatar",
		"created_at",
		"updated_at",
		"deleted_at",
		"stripe_account_id",
	).From("core.users").Where(prefict, params...).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}
	row := runnable.QueryRow(ctx, query, args...)
	scanErr := row.Scan(
		&res.ID,
		&res.UUID,
		&res.StripeID,
		&res.FirstName,
		&res.LastName,
		&res.Phonenumber,
		&res.CountryCode,
		&res.Email,
		&res.Username,
		&res.Type,
		&res.ProfileAvatar,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
		&res.StripeAccountID,
	)
	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return nil, nil
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return nil, scanErr
	}

	return &res, nil
}
