package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/VamaSingapore/vama-api/cmd/test"
	test_util "github.com/VamaSingapore/vama-api/cmd/test/util"
	test_data "github.com/VamaSingapore/vama-api/cmd/test/util/data"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/token"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/gofiber/fiber/v2"
	"github.com/houqp/gtest"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AuthTests struct {
	app      *fiber.App
	db       *pgxpool.Pool
	tokenSvc token.Service
}

func (m *AuthTests) Setup(t *testing.T) {
	fmt.Println("Starting test server for Auth tests")
	app := test.StartTestServer()
	m.app = app.App
	m.db = app.Db

	m.tokenSvc = token.NewService(app.Db, app.Redis, appconfig.Config.Auth.AccessTokenKey)
	test_util.InitializeDbSchemas(m.db)
}

func (m *AuthTests) Teardown(t *testing.T) {
	fmt.Println("Killing test server for Auth tests")
	m.app.Shutdown()
	m.db.Close()
}

func (m *AuthTests) BeforeEach(t *testing.T) {
	test_util.ClearDB(m.db)
}

func (m *AuthTests) AfterEach(t *testing.T) {
}

// Test signing-up a new user /auth/v1/sign-up/sms
func (m *AuthTests) SubTestSignupSMS(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"firstName":   "FirstName",
		"lastName":    "LastName",
		"countryCode": "US",
		"phoneNumber": "9518122389",
		"code":        "12345",
	}

	var asUserID *int
	resp := &response.AuthSuccess{}
	test_util.MakePostRequestAssert200(t, m.app, "/auth/v1/sign-up/sms", params, asUserID, resp)
	expectedSendBirdAccessToken := "fake-access-token-123"
	expectedFirstName := "FirstName"
	expectedLastName := "LastName"
	expectedUserType := "USER"
	expectedPhonenumber := "+19518122389"

	tokenType := "AccessToken"
	tokenClaims, verifyTokenErr := m.tokenSvc.VerifyIDToken(ctx, m.db, resp.Credentials.AccessToken, tokenType)
	if verifyTokenErr != nil {
		vlog.Errorf(ctx, "Error attempting to get token claims: %v", verifyTokenErr)
		t.Fail()
	}

	isValidToken, checkValidErr := test_util.CheckIsValidToken(m.tokenSvc, tokenClaims, resp.Credentials.AccessToken, tokenType)
	if checkValidErr != nil {
		vlog.Errorf(ctx, "Error attempting to verify access token: %v", checkValidErr)
		t.Fail()
	}

	// Token validation assertions
	assert.False(t, tokenClaims.ExpiresAt < time.Now().Unix())
	assert.True(t, isValidToken)

	// Response payload assertions
	assert.Equal(t, expectedSendBirdAccessToken, resp.Credentials.SendBirdAccessToken)
	require.NotNil(t, resp.User)
	assert.Equal(t, expectedFirstName, resp.User.FirstName)
	assert.Equal(t, expectedLastName, resp.User.LastName)
	assert.Equal(t, expectedPhonenumber, resp.User.Phonenumber)
	assert.Equal(t, expectedUserType, resp.User.Type)
	assert.NotNil(t, resp.User.StripeID)
	assert.Nil(t, resp.User.Email)
	assert.Nil(t, resp.User.StripeAccountID)
}

// Test signing-up a net new GOAT user /auth/v1/sign-up/sms/goat
func (m *AuthTests) SubTestSignupSMSNetNewGoat(t *testing.T) {
	// Initializing the payout periods
	test_util.InitializeDbSchemas(m.db)

	// Add Vama User
	ctx := context.Background()
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

	inviteCode := "ABC123"
	test_data.FillBoostrapCreateGoatInviteCodeData(m.db, inviteCode)

	params := map[string]interface{}{
		"firstName":   "FirstName",
		"lastName":    "LastName",
		"countryCode": "US",
		"phoneNumber": "9518172389",
		"code":        "12345",
		"email":       "fake@email.com",
		"username":    "username123",
		"inviteCode":  inviteCode,
	}

	var asUserID *int
	resp := &response.AuthSuccess{}
	test_util.MakePostRequestAssert200(t, m.app, "/auth/v1/sign-up/sms/goat", params, asUserID, resp)

	expectedSendBirdAccessToken := "fake-access-token-123"
	expectedFirstName := "FirstName"
	expectedLastName := "LastName"
	expectedUserType := "GOAT"
	expectedPhonenumber := "+19518172389"
	expectedEmail := "fake@email.com"

	tokenType := "AccessToken"
	tokenClaims, verifyTokenErr := m.tokenSvc.VerifyIDToken(ctx, m.db, resp.Credentials.AccessToken, tokenType)
	if verifyTokenErr != nil {
		vlog.Errorf(ctx, "Error attempting to get token claims: %v", verifyTokenErr)
		t.Fail()
	}

	isValidToken, checkValidErr := test_util.CheckIsValidToken(m.tokenSvc, tokenClaims, resp.Credentials.AccessToken, tokenType)
	if checkValidErr != nil {
		vlog.Errorf(ctx, "Error attempting to verify access token: %s", checkValidErr.Error())
		t.Fail()
	}
	// Token validation assertions
	assert.False(t, tokenClaims.ExpiresAt < time.Now().Unix())
	assert.True(t, isValidToken)

	// Response payload assertions
	assert.Equal(t, expectedSendBirdAccessToken, resp.Credentials.SendBirdAccessToken)
	require.NotNil(t, resp.User)
	assert.Equal(t, expectedFirstName, resp.User.FirstName)
	assert.Equal(t, expectedLastName, resp.User.LastName)
	assert.Equal(t, expectedPhonenumber, resp.User.Phonenumber)
	assert.Equal(t, expectedUserType, resp.User.Type)
	assert.NotNil(t, resp.User.StripeID)
	assert.Equal(t, expectedEmail, *resp.User.Email)
	assert.Nil(t, resp.User.StripeAccountID)
}

// Test upgrading a current user to a creator user /auth/v1/sign-up/sms/goat
func (m *AuthTests) SubTestSignupSMSUpgradeToGoatUser(t *testing.T) {
	// Initializing the payout periods
	test_util.InitializeDbSchemas(m.db)

	// Add Vama User
	ctx := context.Background()
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

	phonenumber := "9518182389"

	// Create a new regular user
	newRegularUserParams := map[string]interface{}{
		"firstName":   "FirstName",
		"lastName":    "LastName",
		"countryCode": "US",
		"phoneNumber": phonenumber,
		"code":        "12345",
	}

	var asUserID *int
	unusedResp := &response.AuthSuccess{}
	test_util.MakePostRequestAssert200(t, m.app, "/auth/v1/sign-up/sms", newRegularUserParams, asUserID, unusedResp)

	// Upgrade the newly created user
	inviteCode := "123ABC"
	test_data.FillBoostrapCreateGoatInviteCodeData(m.db, inviteCode)

	upgradeParams := map[string]interface{}{
		"firstName":   "FirstName",
		"lastName":    "LastName",
		"countryCode": "US",
		"phoneNumber": phonenumber,
		"code":        "12345",
		"email":       "fake12@email.com",
		"username":    "username4321",
		"inviteCode":  inviteCode,
	}

	resp := &response.AuthSuccess{}
	test_util.MakePostRequestAssert200(t, m.app, "/auth/v1/sign-up/sms/goat", upgradeParams, asUserID, resp)

	expectedSendBirdAccessToken := "fake-access-token-123"
	expectedFirstName := "FirstName"
	expectedLastName := "LastName"
	expectedUserType := "GOAT"
	expectedPhonenumber := "+19518182389"
	expectedEmail := "fake12@email.com"

	tokenType := "AccessToken"
	tokenClaims, verifyTokenErr := m.tokenSvc.VerifyIDToken(ctx, m.db, resp.Credentials.AccessToken, tokenType)
	if verifyTokenErr != nil {
		vlog.Errorf(ctx, "Error attempting to get token claims: %s", verifyTokenErr.Error())
		t.Fail()
	}

	isValidToken, checkValidErr := test_util.CheckIsValidToken(m.tokenSvc, tokenClaims, resp.Credentials.AccessToken, tokenType)
	if checkValidErr != nil {
		vlog.Errorf(ctx, "Error attempting to verify access token: %s", checkValidErr.Error())
		t.Fail()
	}

	// Token validation assertions
	assert.False(t, tokenClaims.ExpiresAt < time.Now().Unix())
	assert.True(t, isValidToken)

	// Response payload assertions
	assert.Equal(t, expectedSendBirdAccessToken, resp.Credentials.SendBirdAccessToken)
	require.NotNil(t, resp.User)
	assert.Equal(t, expectedFirstName, resp.User.FirstName)
	assert.Equal(t, expectedLastName, resp.User.LastName)
	assert.Equal(t, expectedPhonenumber, resp.User.Phonenumber)
	assert.Equal(t, expectedUserType, resp.User.Type)
	assert.NotNil(t, resp.User.StripeID)
	assert.Equal(t, expectedEmail, *resp.User.Email)
	assert.Nil(t, resp.User.StripeAccountID)
}

// Test refreshing a user's access token /auth/v1/refresh
func (m *AuthTests) SubTestRefreshToken(t *testing.T) {
	ctx := context.Background()

	// Create a new regular user
	newRegularUserParams := map[string]interface{}{
		"firstName":   "FirstName",
		"lastName":    "LastName",
		"countryCode": "US",
		"phoneNumber": "9528182389",
		"code":        "12345",
	}

	var asUserID *int
	signupResp := &response.AuthSuccess{}
	test_util.MakePostRequestAssert200(t, m.app, "/auth/v1/sign-up/sms", newRegularUserParams, asUserID, signupResp)

	refreshParams := map[string]interface{}{
		"refreshToken": signupResp.Credentials.RefreshToken,
	}

	refreshResp := &response.AuthSuccess{}
	test_util.MakePostRequestAssert200(t, m.app, "/auth/v1/refresh", refreshParams, asUserID, refreshResp)

	tokenType := "AccessToken"
	tokenClaims, verifyTokenErr := m.tokenSvc.VerifyIDToken(ctx, m.db, refreshResp.Credentials.AccessToken, tokenType)
	if verifyTokenErr != nil {
		vlog.Errorf(ctx, "Error attempting to get token claims: %s", verifyTokenErr.Error())
		t.Fail()
	}

	isValidToken, checkValidErr := test_util.CheckIsValidToken(m.tokenSvc, tokenClaims, refreshResp.Credentials.AccessToken, tokenType)
	if checkValidErr != nil {
		vlog.Errorf(ctx, "Error attempting to verify access token: %s", checkValidErr.Error())
		t.Fail()
	}

	assert.False(t, tokenClaims.ExpiresAt < time.Now().Unix())
	assert.True(t, isValidToken)
}

// Test signing-out and revoking a user's access token /auth/v1/sign-out
func (m *AuthTests) SubTestSignout(t *testing.T) {
	ctx := context.Background()

	// Create a new regular user
	newRegularUserParams := map[string]interface{}{
		"firstName":   "FirstName",
		"lastName":    "LastName",
		"countryCode": "US",
		"phoneNumber": "9514182389",
		"code":        "12345",
	}

	var asUserID *int
	signupResp := &response.AuthSuccess{}
	test_util.MakePostRequestAssert200(t, m.app, "/auth/v1/sign-up/sms", newRegularUserParams, asUserID, signupResp)

	tokenType := "RefreshToken"
	tokenClaims, verifyTokenErr := m.tokenSvc.VerifyIDToken(ctx, m.db, signupResp.Credentials.RefreshToken, tokenType)
	if verifyTokenErr != nil {
		vlog.Errorf(ctx, "Error attempting to get token claims: %s", verifyTokenErr.Error())
		t.Fail()
	}

	isValidToken, checkValidErr := test_util.CheckIsValidToken(m.tokenSvc, tokenClaims, signupResp.Credentials.RefreshToken, tokenType)
	if checkValidErr != nil {
		vlog.Errorf(ctx, "Error attempting to verify access token: %s", checkValidErr.Error())
		t.Fail()
	}

	// Ensure we start with a valid access token
	assert.False(t, tokenClaims.ExpiresAt < time.Now().Unix())
	assert.True(t, isValidToken)

	// Sign-out and revoke access token
	signoutParams := map[string]interface{}{
		"refreshToken": signupResp.Credentials.RefreshToken,
	}

	unusedResp := &response.AuthSuccess{}
	test_util.MakePostRequestAssert200(t, m.app, "/auth/v1/sign-out", signoutParams, asUserID, unusedResp)

	signoutTokenType := "RefreshToken"
	signoutIsValidToken, signoutCheckValidErr := test_util.CheckIsValidToken(m.tokenSvc, tokenClaims, signupResp.Credentials.RefreshToken, signoutTokenType)
	if signoutCheckValidErr != nil {
		vlog.Errorf(ctx, "Error attempting to verify refresh token: %s", signoutCheckValidErr.Error())
		t.Fail()
	}

	assert.False(t, signoutIsValidToken)
}

func TestAuth(t *testing.T) {
	gtest.RunSubTests(t, &AuthTests{})
}
