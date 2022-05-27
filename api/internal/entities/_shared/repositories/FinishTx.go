package repositories

import (
	"context"
	"errors"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

/*
This function should be used as follows:

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when {doing whatever}.",
			fmt.Sprintf("Could not begin transaction: %v", txErr),
		)
	}
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)
	... (code using transaction etc etc)
	commit = true
	return {whatever}

Only set commit = true in the successful return
cases. The transaction will rollback in the error
cases by default.
*/

func FinishTx(ctx context.Context, tx pgx.Tx, commit *bool) error {
	var txErr error
	if commit == nil {
		tx.Rollback(ctx)
		commitErr := errors.New("commit was nil")
		vlog.Errorf(ctx, commitErr.Error())
		return commitErr
	}
	if *commit {
		txErr = tx.Commit(ctx)
	} else {
		txErr = tx.Rollback(ctx)
	}
	if txErr != nil {
		vlog.Errorf(ctx, "Error finishing transaction. commit = %v: %v", *commit, txErr)
	}
	return txErr
}
