package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) InsertPendingContact(ctx context.Context, userID int, phone request.Phone, firstName, lastName string) error {
	query, args, squirrelErr := squirrel.Insert("core.pending_contacts").
		Columns(
			"user_id",
			"phone_number",
			"country_code",
			"first_name",
			"last_name",
		).
		Values(
			userID,
			phone.Number,
			phone.CountryCode,
			firstName,
			lastName,
		).
		Suffix("ON CONFLICT DO NOTHING").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := s.MasterNode().Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
