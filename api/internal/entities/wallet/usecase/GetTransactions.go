package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getTransactionsDefaultErr = "Something went wrong getting transactions history."

func (svc *usecase) GetTransactions(ctx context.Context, req request.GetTransactions) (*response.GetTransactions, error) {
	currentUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getTransactionsDefaultErr,
			fmt.Sprintf("Error occurred when retrieving user from db. Err: %v", getCurrUserErr),
		)
	} else if currentUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusInternalServerError,
			getTransactionsDefaultErr,
			"Current user nil for GetTransactions",
		)
	}

	transactions, getTransactionsErr := svc.repo.GetTransactionHistory(ctx, currentUser.ID, req.CursorID, uint64(req.Limit))
	if getTransactionsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getTransactionsDefaultErr,
			fmt.Sprintf("Error getting transaction history from ledger: %v", getTransactionsErr),
		)
	}

	return &response.GetTransactions{Transactions: transactions}, nil
}
