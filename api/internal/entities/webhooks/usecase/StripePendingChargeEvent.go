package service

import (
	"context"
	"net/http"

	walletrepo "github.com/VamaSingapore/vama-api/internal/entities/wallet/repositories"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/stripe/stripe-go/v72"
)

func (svc *usecase) StripePendingChargeEvent(ctx context.Context, event stripe.Event) error {
	chargeID := event.Data.Object["id"].(string)

	getChargeParams := &stripe.ChargeParams{}
	getChargeParams.AddExpand("payment_intent")
	charge, getChargeErr := svc.stripeClient.Charges.Get(chargeID, getChargeParams)
	if getChargeErr != nil {
		vlog.Errorf(ctx, "Error getting charge ID %s: %s", chargeID, getChargeErr.Error())
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong getting pending charge",
		)
	}

	if event.Type == "charge.expired" {
		deleteChargeErr := walletrepo.DeletePendingPayment(ctx, svc.repo.MasterNode(), charge.PaymentIntent.ID)
		if deleteChargeErr != nil {
			vlog.Errorf(ctx, "Error deleting pending payment %s: %s", charge.PaymentIntent.ID, deleteChargeErr.Error())
			return httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				"Something went wrong deleting pending charge",
			)
		}
	}
	return nil
}
