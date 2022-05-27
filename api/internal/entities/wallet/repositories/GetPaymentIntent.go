package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetPaymentIntent(ctx context.Context, customerUserID int, providerUserID int) (string, error) {
	query, args, squirrelErr := squirrel.Select(
		"payment_intent_id",
	).
		From("wallet.pending_transactions").
		Where(`id = (
			SELECT MAX(id)
			FROM wallet.pending_transactions
			WHERE customer_user_id = ? AND provider_user_id = ?
		)`, customerUserID, providerUserID).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if squirrelErr != nil {
		return "", squirrelErr
	}

	row := s.MasterNode().QueryRow(ctx, query, args...)

	var paymentIntentID string
	scanErr := row.Scan(
		&paymentIntentID,
	)
	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return "", constants.ErrNotFound
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return "", scanErr
	}
	return paymentIntentID, nil
}
