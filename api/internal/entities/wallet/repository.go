package wallet

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stripe/stripe-go/v72"
)

type Repository interface {
	baserepo.BaseRepository
	MasterNode() *pgxpool.Pool
	GetGoatChatPrice(ctx context.Context, goatUserID int) (*response.GetGoatChatPrice, error)
	GetPaymentIntent(ctx context.Context, customerUserID int, providerUserID int) (string, error)
	UpsertBillingAddress(ctx context.Context, tx pgx.Tx, userID int, req request.BillingAddress) error
	UpsertBank(ctx context.Context, tx pgx.Tx, userID int, req request.UpsertBank) error
	GetUserBankInfo(ctx context.Context, userID int) ([]response.PaymentMethod, error)
	GetPayoutPeriods(ctx context.Context, req request.GetPayoutPeriods) ([]response.PayoutPeriod, error)
	GetPayPeriodPerTS(ctx context.Context, ts int64) (*response.PayoutPeriod, error)
	ListPayoutHistory(ctx context.Context, req request.ListPayoutHistory, limit int64) ([]response.PayoutHistoryDatum, error)

	// ledger methods
	GetBalance(ctx context.Context, providerUserID int, currency string) (*int64, error)
	UpsertBalance(ctx context.Context, runnable utils.Runnable, ledgerEntry LedgerEntry) error

	InsertPendingPayment(ctx context.Context, exec utils.Executable, paymentIntent stripe.PaymentIntent, customerID int, providerID int) error
	DeletePendingPayment(ctx context.Context, exec utils.Executable, paymentIntentID string) error
	GetPendingTransactionHistory(ctx context.Context, userID int, lastTransactionID int64, limit uint64) ([]response.TransactionItem, error)

	UpdateBalancePayout(ctx context.Context, tx pgx.Tx, providerID int, amountPaid int64, payPeriodID int) error
	ListUnpaidProviders(ctx context.Context) ([]response.ProviderPaymentInfo, error)
	InsertLedgerTransaction(ctx context.Context, runnable utils.Runnable, ledgerEntry LedgerEntry) error
	GetTransactionHistory(ctx context.Context, userID int, lastTransactionID int64, limit uint64) ([]response.TransactionItem, error)
	GetPendingBalanceNotifications(ctx context.Context, runnable utils.Runnable) ([]PendingBalanceNotification, error)
}
