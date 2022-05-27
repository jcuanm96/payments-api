package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) DeletePendingPayment(ctx context.Context, exec utils.Executable, paymentIntentID string) error {
	return DeletePendingPayment(ctx, exec, paymentIntentID)
}

func DeletePendingPayment(ctx context.Context, exec utils.Executable, paymentIntentID string) error {
	query, args, squirrelErr := squirrel.Delete("wallet.pending_transactions").
		Where("payment_intent_id = ?", paymentIntentID).
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
