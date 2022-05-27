package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/stripe/stripe-go/v72"
)

func (svc *usecase) GetNewBalanceTransaction(ctx context.Context, extractedStripeEvent wallet.StripeTxnEvent) error {
	// Trial invoices have Amount 0 and do not generate a balance transaction; do nothing.
	if extractedStripeEvent.Metadata != nil {
		isTrial := extractedStripeEvent.Metadata.IsTrial
		if isTrial != nil && *isTrial {
			return nil
		}
	}

	var ledgerEntry *wallet.LedgerEntry
	var paymentIntentID string
	hasMore := true
	for hasMore {
		params := &stripe.BalanceTransactionListParams{
			CreatedRange: &stripe.RangeQueryParams{
				GreaterThan: extractedStripeEvent.CreatedAt - 60, // Using a buffer of 60 seconds
			},
		}

		// https://stripe.com/docs/api/balance_transactions/list#balance_transaction_list-limit
		limit := int64(100)
		params.Filters.AddFilter("limit", "", strconv.FormatInt(limit, 10))
		// include full payment source object instead of just its ID
		// to avoid making another request for each transaction
		params.AddExpand("data.source")
		params.AddExpand("data.source.payment_intent")
		params.AddExpand("data.source.invoice")

		i := svc.stripeClient.BalanceTransaction.List(params)
		hasMore = i.BalanceTransactionList().HasMore

		for i.Next() {
			transaction := i.BalanceTransaction()
			shouldBreak := false

			payPeriod, getPayPeriodErr := svc.repo.GetPayPeriodPerTS(ctx, transaction.Created)
			if getPayPeriodErr != nil {
				payPeriodErrMsg := fmt.Sprintf("Error getting pay period for transaction id The %s and ts %d. Err: %s", transaction.ID, transaction.Created, getPayPeriodErr)
				vlog.Errorf(ctx, payPeriodErrMsg)
				telegram.TelegramClient.SendMessage(payPeriodErrMsg)
				payPeriod = &response.PayoutPeriod{
					ID: -1,
				}
			} else if payPeriod == nil {
				payPeriodErrMsg := fmt.Sprintf("Error pay period came back as nil for id %s and timestamp %d.", transaction.ID, transaction.Created)
				vlog.Errorf(ctx, payPeriodErrMsg)
				telegram.TelegramClient.SendMessage(payPeriodErrMsg)
				payPeriod = &response.PayoutPeriod{
					ID: -1,
				}
			}
			switch transaction.Source.Type {
			case stripe.BalanceTransactionSourceTypeCharge:
				if transaction.Source.Charge.PaymentIntent == nil {
					vlog.Errorf(ctx, "payment intent for event the following charge payload came back as nil: %v", transaction.Source.Charge)
					continue
				}

				if extractedStripeEvent.ID != transaction.Source.Charge.ID &&
					(transaction.Source.Charge.Invoice != nil &&
						extractedStripeEvent.ID != transaction.Source.Charge.Invoice.ID) {
					continue
				}

				var customerUserID int
				var providerUserID int

				// Non-nil invoice means this is a subscription,
				// nil invoice means this is a creator chat.
				if transaction.Source.Charge.Invoice != nil {
					if extractedStripeEvent.Metadata == nil {
						errorMsg := fmt.Sprintf("metadata for Subscription invoice was nil: %s", extractedStripeEvent.ID)
						return errors.New(errorMsg)
					}

					_customerUserID := extractedStripeEvent.Metadata.CustomerUserID
					if _customerUserID == 0 {
						return errors.New("customerUserID was 0 in extracted subscription metadata")
					}
					customerUserID = _customerUserID

					_providerUserID := extractedStripeEvent.Metadata.ProviderUserID
					if _providerUserID == 0 {
						return errors.New("providerUserID was 0 in extracted subscription metadata")
					}
					providerUserID = _providerUserID
				} else {
					metadataValues, getIDsFromMetadataErr := wallet.GetUserIDsFromStripeMetadata(ctx, transaction.Source.Charge.PaymentIntent.Metadata)
					if getIDsFromMetadataErr != nil {
						errorMsg := fmt.Sprintf("Something went wrong getting user IDs from metadata in payment intent %s: %v", transaction.Source.Charge.PaymentIntent.ID, getIDsFromMetadataErr)
						return errors.New(errorMsg)
					}

					providerUserID = metadataValues.ProviderUserID
					customerUserID = metadataValues.CustomerUserID
					// Only record the paymentIntentID to delete if this was a manually captured pending charge,
					// as automatic captures aren't recorded as pending in our DB.
					if transaction.Source.Charge.PaymentIntent.CaptureMethod == stripe.PaymentIntentCaptureMethodManual {
						paymentIntentID = transaction.Source.Charge.PaymentIntent.ID
					}
				}

				totalFeesRatio := constants.TOTAL_FEES_RATIO
				const noFeesUserID = 292
				if appconfig.Config.Gcloud.Project == "vama-prod" && providerUserID == noFeesUserID {
					totalFeesRatio = 0
				}

				totalFeesExact := float64(transaction.Amount) * totalFeesRatio
				vamaFee := int64(math.Ceil(totalFeesExact)) - transaction.Fee

				ledgerEntry = &wallet.LedgerEntry{
					ProviderUserID:      providerUserID,
					CustomerUserID:      customerUserID,
					StripeTransactionID: transaction.ID,
					SourceType:          string(transaction.Source.Type),
					Amount:              transaction.Amount,
					VamaFee:             vamaFee,
					StripeFee:           transaction.Fee,
					CreatedTS:           transaction.Created,
					BalanceDelta:        transaction.Amount - int64(totalFeesExact),
					Currency:            string(transaction.Currency),
					PayPeriodID:         payPeriod.ID,
				}
				shouldBreak = true
			}
			if shouldBreak {
				break
			}
		}
	}

	if ledgerEntry == nil {
		errorMsg := fmt.Sprintf("failed to assign ledger entry a value for Stripe event: %v", extractedStripeEvent)
		return errors.New(errorMsg)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		errorMsg := fmt.Sprintf("Could not begin transaction for GetNewBalanceTransaction. Error: %v", txErr)
		return errors.New(errorMsg)
	}

	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	updateBalanceErr := svc.repo.UpsertBalance(ctx, tx, *ledgerEntry)
	if updateBalanceErr != nil {
		return updateBalanceErr
	}

	if ledgerEntry.SourceType == string(stripe.BalanceTransactionSourceTypeCharge) &&
		paymentIntentID != "" {
		deletePaymentErr := svc.repo.DeletePendingPayment(ctx, tx, paymentIntentID)
		if deletePaymentErr != nil {
			errorMsg := fmt.Sprintf("Error deleting pending payment %s: %v", paymentIntentID, deletePaymentErr)
			vlog.Errorf(ctx, errorMsg)
			telegram.TelegramClient.SendMessage(errorMsg)
		}
	}

	insertTransactionErr := svc.repo.InsertLedgerTransaction(ctx, tx, *ledgerEntry)
	if insertTransactionErr != nil {
		errorMsg := fmt.Sprintf("Error inserting new transaction: %v", insertTransactionErr)
		return errors.New(errorMsg)
	}

	commit = true
	return nil
}
