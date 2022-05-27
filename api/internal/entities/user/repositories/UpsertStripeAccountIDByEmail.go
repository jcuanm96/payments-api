package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) UpsertStripeAccountIDByEmail(ctx context.Context, email string, stripeAccountID string) error {
	query, args, squirrelErr := squirrel.Update("core.users").
		Set("stripe_account_id", stripeAccountID).
		Where("email=?", email).
		PlaceholderFormat(squirrel.Dollar).ToSql()

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
