package wallet

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/stripe/stripe-go/v72"
)

type Usecase interface {
	UpsertBank(ctx context.Context, req request.UpsertBank) error
	SaveDefaultPaymentMethod(ctx context.Context, req request.DefaultPaymentMethod) (*response.DefaultPaymentMethod, error)
	GetMyPaymentMethods(ctx context.Context) (*response.GetMyPaymentMethods, error)

	MakePaymentIntent(ctx context.Context, req request.MakePaymentIntent) error
	MakeChatPaymentIntent(ctx context.Context, req request.MakeChatPaymentIntent) error
	ConfirmPaymentIntent(ctx context.Context, req request.ConfirmPaymentIntent) (*response.ConfirmPaymentIntent, error)
	RefundSubscription(ctx context.Context, paymentIntentID string) (*stripe.Refund, error)

	GetGoatChatPrice(ctx context.Context, req request.GetGoatChatPrice) (*response.GetGoatChatPrice, error)
	UpsertGoatChatsPrice(ctx context.Context, req request.UpsertGoatChatsPrice) (*response.UpsertGoatChatsPrice, error)

	GetPayoutPeriods(ctx context.Context, req request.GetPayoutPeriods) (*response.GetPayoutPeriods, error)
	ListPayoutHistory(ctx context.Context, req request.ListPayoutHistory) (*response.ListPayoutHistory, error)
	GetPayPeriod(ctx context.Context, req request.GetPayPeriod) (*response.PayoutPeriod, error)

	// ledger methods
	GetBalance(ctx context.Context, req request.GetBalance) (*response.GetBalance, error)
	GetNewBalanceTransaction(ctx context.Context, extractedStripeEvent StripeTxnEvent) error
	GetTransactions(ctx context.Context, req request.GetTransactions) (*response.GetTransactions, error)
	GetPendingTransactions(ctx context.Context, req request.GetTransactions) (*response.GetTransactions, error)
	MarkProviderAsPaid(ctx context.Context, req request.MarkProviderAsPaid) error
	ListUnpaidProviders(ctx context.Context) (*response.ListUnpaidProviders, error)
	UpsertBalance(ctx context.Context, runnable utils.Runnable, ledgerEntry LedgerEntry) error
	InsertLedgerTransaction(ctx context.Context, runnable utils.Runnable, ledgerEntry LedgerEntry) error

	SendPendingBalanceNotifications(ctx context.Context) error
}
