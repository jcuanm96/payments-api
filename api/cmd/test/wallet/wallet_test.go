package test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/VamaSingapore/vama-api/cmd/test"
	"github.com/VamaSingapore/vama-api/cmd/test/mocks"
	test_util "github.com/VamaSingapore/vama-api/cmd/test/util"
	test_data "github.com/VamaSingapore/vama-api/cmd/test/util/data"
	"github.com/VamaSingapore/vama-api/internal/controller"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/auth"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/gofiber/fiber/v2"
	"github.com/houqp/gtest"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	walletrepo "github.com/VamaSingapore/vama-api/internal/entities/wallet/repositories"
	walletusecase "github.com/VamaSingapore/vama-api/internal/entities/wallet/usecase"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

type WalletTests struct {
	app *fiber.App
	ctr *controller.Ctr
	db  *pgxpool.Pool
}

func (m *WalletTests) Setup(t *testing.T) {
	fmt.Println("Starting test server for Wallet tests")
	app := test.StartTestServer()
	m.app = app.App
	m.db = app.Db
	m.ctr = app.Ctr
	test_util.InitializeDbSchemas(m.db)
}

func (m *WalletTests) Teardown(t *testing.T) {
	fmt.Println("Killing test server for Wallet tests")
	m.app.Shutdown()
	m.db.Close()
}

func (m *WalletTests) BeforeEach(t *testing.T) {
	test_util.ClearDB(m.db)
}

func (m *WalletTests) AfterEach(t *testing.T) {
}

// Test getting a users balance /wallet/balance/me
// returns a 200 w/ a valid balance
func (m *WalletTests) SubTestGetBalance(t *testing.T) {
	test_data.FillBootstrapUserData(m.db)

	asUserID := 1
	balanceDeltaInSmallestDenom := int64(100)

	test_data.FillBootstrapWalletBalanceData(m.db, asUserID, balanceDeltaInSmallestDenom)

	params := map[string]string{
		"currency": "usd",
	}
	resp := &response.GetBalance{}

	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/wallet/balance/me", params, asUserID, resp)

	assert.Equal(t, "usd", resp.Currency)
	assert.Equal(t, balanceDeltaInSmallestDenom, resp.Amount)
}

// Test getting a users transaction history /wallet/transactions/me
// returns a 200 and confirms correct ledger items in ascending order
func (m *WalletTests) SubTestGetUserTransactionHistory(t *testing.T) {
	test_data.FillBootstrapUserData(m.db)

	asProviderUserID := 1
	customerUserID := 2
	test_data.FillBootstrapTransactionItems(m.db, asProviderUserID, customerUserID)
	test_data.FillBootstrapTransactionItems(m.db, asProviderUserID, customerUserID)

	params := map[string]string{}
	resp := &response.GetTransactions{}

	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/wallet/transactions/me", params, asProviderUserID, resp)

	require.Equal(t, 2, len(resp.Transactions))
	assert.True(t, resp.Transactions[0].CreatedAt <= resp.Transactions[1].CreatedAt)
	assert.Equal(t, int64(50), resp.Transactions[0].Fees)
	assert.Equal(t, "usd", resp.Transactions[0].Currency)
	assert.Equal(t, constants.TXN_HISTORY_OUTGOING, resp.Transactions[0].Type)
	require.NotNil(t, resp.Transactions[0].User)
	assert.Equal(t, 2, resp.Transactions[0].User.ID)
}

// Test that the user's balance gets updated after a successful payment intent confirmation
// returns a 200 and confirms correct ledger items in ascending order w/ correct updated balance
func (m *WalletTests) SubTestGetNewBalanceTransaction(t *testing.T) {
	// Re-initialzing the tests so that we can populate the wallet.payout_periods
	ctx := context.Background()
	test_util.InitializeDbSchemas(m.db)

	// Add Vama User
	tx, txErr := m.db.Begin(ctx)
	if txErr != nil {
		vlog.Fatalf(ctx, "Error starting transaction in SubTestSignupSMSNetNewGoat: %s", txErr.Error())
	}
	insertVamaUserErr := test_data.FillBootstrapVamaUser(tx)
	if insertVamaUserErr != nil {
		vlog.Fatalf(ctx, "Error inserting Vama user: %s", insertVamaUserErr.Error())
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		vlog.Fatalf(ctx, "Error commiting: %s", commitErr.Error())
	}

	asProviderUserID := 2
	asCustomerUserID := 1

	signupResp := test_data.FillBootstrapCustomerUserSignupData(t, m.app)

	require.NotNil(t, signupResp.User)
	require.Equal(t, asCustomerUserID, signupResp.User.ID)

	// Create a new creator user
	inviteCode := "HIJ123"
	test_data.FillBoostrapCreateGoatInviteCodeData(m.db, inviteCode)
	goatSignupResp := test_data.FillBootstrapGoatUserSignupData(t, m.app, inviteCode)

	require.NotNil(t, goatSignupResp.User)
	require.Equal(t, asProviderUserID, goatSignupResp.User.ID)

	// Save default payment method
	thisYear := time.Now().Year()
	expirationYear := thisYear + 8 // There's a limit on how large the year can be, so randomly chose offset
	savePaymentMethodParams := map[string]interface{}{
		"number":   "4242424242424242",
		"expMonth": "01",
		"expYear":  strconv.Itoa(expirationYear),
		"cvc":      "123",
	}
	savePaymentMethodResp := &response.DefaultPaymentMethod{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/wallet/payment-method/credit-card/me", savePaymentMethodParams, &asCustomerUserID, savePaymentMethodResp)

	// Make payment intent
	makeChatPaymentIntentParams := map[string]interface{}{
		"providerUserID": asProviderUserID,
	}
	makeChatPaymentIntentResp := &struct{}{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/wallet/payment-intents/goat/chat", makeChatPaymentIntentParams, &asCustomerUserID, makeChatPaymentIntentResp)

	// Confirm payment intent
	confirmPaymentIntentParams := map[string]interface{}{
		"customerUserID": asCustomerUserID,
	}
	confirmPaymentIntentResp := &response.ConfirmPaymentIntent{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/wallet/payment-intents/goat/chat/confirm", confirmPaymentIntentParams, &asProviderUserID, confirmPaymentIntentResp)

	// GetNewBalanceTransaction invocation
	var useruc user.Usecase
	var authuc auth.Usecase
	pushuc := mocks.NewMockPush()
	walletuc := walletusecase.New(walletrepo.New(m.db), useruc, &m.ctr.Stripe, authuc, pushuc)
	stripeTxEvent := wallet.StripeTxnEvent{
		ID:        confirmPaymentIntentResp.ChargeID,
		CreatedAt: time.Now().Unix() - 120, // Adding buffer time of 2 minutes
		Metadata: &wallet.StripeEventMetadata{
			ProviderUserID: asProviderUserID,
			CustomerUserID: asCustomerUserID,
		},
	}
	getNewBalanceTxErr := walletuc.GetNewBalanceTransaction(context.Background(), stripeTxEvent)

	if getNewBalanceTxErr != nil {
		vlog.Fatalf(ctx, "Error committing transaction: %s", getNewBalanceTxErr.Error())
	}

	getBalanceParams := map[string]string{
		"currency": "usd",
	}
	getBalanceResp := &response.GetBalance{}

	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/wallet/balance/me", getBalanceParams, asProviderUserID, getBalanceResp)

	expectedFinalBalanceCredit := 285
	assert.Equal(t, int64(constants.GOAT_SIGNUP_CREDIT_AMOUNT_USD_IN_SMALLEST_DENOM+expectedFinalBalanceCredit), getBalanceResp.Amount)
}

func TestWallet(t *testing.T) {
	gtest.RunSubTests(t, &WalletTests{})
}
