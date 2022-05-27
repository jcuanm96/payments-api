package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/stripe/stripe-go/v72"
)

func (s *repository) InsertPendingPayment(ctx context.Context, exec utils.Executable, paymentIntent stripe.PaymentIntent, customerID int, providerID int) error {
	query, args, squirrelErr := squirrel.Insert("wallet.pending_transactions").
		Columns(
			"provider_user_id",
			"customer_user_id",
			"payment_intent_id",
			"stripe_created_ts",
			"amount",
			"currency",
			"version",
		).
		Values(
			providerID,
			customerID,
			paymentIntent.ID,
			paymentIntent.Created,
			paymentIntent.Amount,
			paymentIntent.Currency,
			constants.CURR_PAYMENTS_VERSION,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := exec.Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
