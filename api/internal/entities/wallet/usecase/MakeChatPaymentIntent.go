package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) MakeChatPaymentIntent(ctx context.Context, req request.MakeChatPaymentIntent) error {
	// Get creator chat price from DB
	goatChatPrice, getGoatChatPriceErr := svc.repo.GetGoatChatPrice(ctx, req.ProviderUserID)

	if getGoatChatPriceErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			makePaymentIntentDefaultErr,
			fmt.Sprintf("Error getting creator chat price from db for provider %d. Err: %v", req.ProviderUserID, getGoatChatPriceErr),
		)
	} else if goatChatPrice.PriceInSmallestDenom < constants.MIN_GOAT_CHAT_PRICE_USD_IN_SMALLEST_DENOM {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			makePaymentIntentDefaultErr,
			fmt.Sprintf("Error the price for a creator chat for provider %d is too low.", req.ProviderUserID),
		)
	}

	makePaymentIntentReq := request.MakePaymentIntent{
		ProviderUserID:        req.ProviderUserID,
		AmountInSmallestDenom: goatChatPrice.PriceInSmallestDenom,
		Currency:              goatChatPrice.Currency,
		AutoCapture:           false,
	}
	makePaymentIntentHTTPErr := svc.MakePaymentIntent(ctx, makePaymentIntentReq)
	if makePaymentIntentHTTPErr != nil {
		return makePaymentIntentHTTPErr
	}

	return nil
}
