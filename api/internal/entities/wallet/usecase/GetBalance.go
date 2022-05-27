package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const getBalanceErr = "There was a problem getting your balance."

func (svc *usecase) GetBalance(ctx context.Context, req request.GetBalance) (*response.GetBalance, error) {
	currentUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getBalanceErr,
			fmt.Sprintf("Something went wrong when getting the user to retrieve balance. Err: %v", getCurrUserErr),
		)
	} else if currentUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			getBalanceErr,
			"Current user is nil for GetBalance.",
		)
	} else if currentUser.Email == nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You must set your email before being able to get your balance.",
			"You must set your email before being able to get your balance.",
		)
	}

	var wg sync.WaitGroup

	var balance *int64
	var balanceErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		balance, balanceErr = svc.repo.GetBalance(ctx, currentUser.ID, req.Currency)
	}()

	wg.Add(1)
	var banks []response.PaymentMethod
	var getBanksErr error
	go func() {
		defer wg.Done()
		banks, getBanksErr = svc.repo.GetUserBankInfo(ctx, currentUser.ID)
	}()

	wg.Wait()

	if balanceErr != nil {
		vlog.Errorf(ctx, "Error getting balance: %s", balanceErr.Error())
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getBalanceErr,
			fmt.Sprintf("Something went wrong getting user balance: %v", balanceErr),
		)
	}
	var balanceAmount int64
	if balance == nil {
		balanceAmount = 0
	} else {
		balanceAmount = *balance
	}

	res := response.GetBalance{
		Currency: req.Currency,
		Amount:   balanceAmount,
	}

	if getBanksErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getBalanceErr,
			fmt.Sprintf("Something went wrong getting banks: %v", getBanksErr),
		)
	}

	if len(banks) > 0 {
		res.Bank = banks[0].Bank
	}

	// Need to specify form as string: https://pkg.go.dev/time#example-Date
	const initialPayoutDate = "2022-02-09T12:00:00Z"
	payoutDate, parseErr := time.Parse(time.RFC3339, initialPayoutDate)
	if parseErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			getBalanceErr,
			fmt.Sprintf("Something went wrong parsing time string in GetBalance %s: %v", initialPayoutDate, parseErr),
		)
	}

	// Get the next payout date (every 2 weeks)
	for payoutDate.Before(time.Now()) {
		const years = 0
		const months = 0
		const days = 14
		payoutDate = payoutDate.AddDate(years, months, days)
	}

	res.NextPayout = payoutDate

	return &res, nil
}
