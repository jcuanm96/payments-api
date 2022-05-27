package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getPendingTransactionsDefaultErr = "Something went wrong getting pending transactions history."

func (svc *usecase) GetPendingTransactions(ctx context.Context, req request.GetTransactions) (*response.GetTransactions, error) {
	currentUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getPendingTransactionsDefaultErr,
			fmt.Sprintf("Error occurred when retrieving user from db. Err: %v", getCurrUserErr),
		)
	} else if currentUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusInternalServerError,
			getPendingTransactionsDefaultErr,
			"Current user nil for GetPendingTransactions",
		)
	}

	transactions, getTransactionsErr := svc.repo.GetPendingTransactionHistory(ctx, currentUser.ID, req.CursorID, uint64(req.Limit))
	if getTransactionsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getPendingTransactionsDefaultErr,
			fmt.Sprintf("Error getting pending transaction history from ledger: %v", getTransactionsErr),
		)
	}

	return &response.GetTransactions{Transactions: transactions}, nil
}
