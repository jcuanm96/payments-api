package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const upsertGoatChatPriceErr = "Something went wrong when updating chat price."

func (svc *usecase) UpsertGoatChatsPrice(ctx context.Context, req request.UpsertGoatChatsPrice) (*response.UpsertGoatChatsPrice, error) {
	currentUser, userErr := svc.user.GetCurrentUser(ctx)

	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			upsertGoatChatPriceErr,
			fmt.Sprintf("Error occurred when getting the user to update creator chats price. Err: %v", userErr),
		)
	} else if currentUser.Type != "GOAT" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can charge for chats. Interested?  Sign up to be a creator now!",
			"Only creators can charge for chats.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			upsertGoatChatPriceErr,
			fmt.Sprintf("Could not begin transaction. Error: %v", txErr),
		)
	}

	// Defer a rollback unless we set commit to true
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	priceErr := svc.user.UpsertGoatChatsPrice(ctx, tx, req.PriceInSmallestDenom, req.Currency, currentUser.ID)
	if priceErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			upsertGoatChatPriceErr,
			fmt.Sprintf("Error occurred when updating creator chats price with amount %d and currency %s from db for user %d. Err: %v", req.PriceInSmallestDenom, req.Currency, currentUser.ID, priceErr),
		)
	}

	res := response.UpsertGoatChatsPrice{}
	commit = true
	return &res, nil
}
