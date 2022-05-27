package test_util

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	test_util "github.com/VamaSingapore/vama-api/cmd/test/util"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	contactrepo "github.com/VamaSingapore/vama-api/internal/entities/contact/repositories"
	feedrepo "github.com/VamaSingapore/vama-api/internal/entities/feed/repositories"
	followrepo "github.com/VamaSingapore/vama-api/internal/entities/follow/repositories"
	userrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	walletrepo "github.com/VamaSingapore/vama-api/internal/entities/wallet/repositories"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Create a new user in the db that contains valid Stripe credentials
func FillBootstrapCustomerUserSignupData(t *testing.T, app *fiber.App) *response.AuthSuccess {
	// Create a new customer user
	newRegularUserParams := map[string]interface{}{
		"firstName":   "FirstName",
		"lastName":    "LastName",
		"countryCode": "US",
		"phoneNumber": "9318182389",
		"code":        "12345",
	}

	signupResp := &response.AuthSuccess{}
	test_util.MakePostRequestAssert200(t, app, "/auth/v1/sign-up/sms", newRegularUserParams, nil, signupResp)

	// Update user email
	updateEmailParams := map[string]interface{}{
		"email": "fake@email96.com",
	}

	updateEmailResp := &response.User{}
	test_util.MakePatchRequestAssert200(t, app, "/api/v1/users/me", updateEmailParams, nil, updateEmailResp)
	return signupResp
}

// Create a new creator user in the db that contains valid Stripe credentials
func FillBootstrapGoatUserSignupData(t *testing.T, app *fiber.App, inviteCode string) *response.AuthSuccess {
	params := map[string]interface{}{
		"firstName":   "FirstName",
		"lastName":    "LastName",
		"countryCode": "US",
		"phoneNumber": "8168172389",
		"code":        "12345",
		"email":       "fake@goatemail20.com",
		"username":    "bitcoinMcQuackers123",
		"inviteCode":  inviteCode,
	}

	var asUserID *int
	goatSignupResp := &response.AuthSuccess{}
	test_util.MakePostRequestAssert200(t, app, "/auth/v1/sign-up/sms/goat", params, asUserID, goatSignupResp)
	return goatSignupResp
}

func FillBootstrapUserData(db *pgxpool.Pool) {
	ctx := context.Background()
	tx, txErr := db.Begin(ctx)
	if txErr != nil {
		vlog.Fatalf(ctx, "Error starting transaction: %s", txErr.Error())
	}

	email1 := "user1@testemail.com"
	stripeId1 := "cus_KJSVHC5V448NO"
	user1 := response.User{
		FirstName:   "user1first",
		LastName:    "user1last",
		Phonenumber: "1234567890",
		CountryCode: "US",
		Email:       &email1,
		Username:    "user1-username",
		Type:        "GOAT",
		StripeID:    &stripeId1,
	}
	_, userDbErr := userrepo.UpsertUserDB(ctx, tx, &user1)
	if userDbErr != nil {
		vlog.Fatalf(ctx, "Error upserting user: %v", userDbErr)
	}

	email2 := "user2@testemail.com"
	stripeId2 := "cus_Kz2mW10anmNin0"
	user2 := response.User{
		FirstName:   "user2first",
		LastName:    "user2last",
		Phonenumber: "9876543210",
		CountryCode: "US",
		Email:       &email2,
		Username:    "user2-username",
		Type:        "USER",
		StripeID:    &stripeId2,
	}
	_, userDbErr2 := userrepo.UpsertUserDB(ctx, tx, &user2)
	if userDbErr2 != nil {
		vlog.Fatalf(ctx, "Error upserting user: %v", userDbErr2)
	}

	email3 := "user3@testemail.com"
	stripeId3 := "cus_JJJVHC5V448NO"
	user3 := response.User{
		FirstName:   "user3first",
		LastName:    "user3last",
		Phonenumber: "9876543310",
		CountryCode: "US",
		Email:       &email3,
		Username:    "user3-username",
		Type:        "USER",
		StripeID:    &stripeId3,
	}
	_, userDbErr3 := userrepo.UpsertUserDB(ctx, tx, &user3)
	if userDbErr3 != nil {
		vlog.Fatalf(ctx, "Error upserting user: %s", userDbErr3.Error())
	}

	insertVamaUserErr := FillBootstrapVamaUser(tx)
	if insertVamaUserErr != nil {
		vlog.Fatalf(ctx, "Error inserting Vama user: %s", insertVamaUserErr.Error())
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		vlog.Fatalf(ctx, "Error committing transaction: %s", commitErr.Error())
	}
}

func FillBootstrapFeedPostData(db *pgxpool.Pool) {
	ctx := context.Background()
	postReq := request.MakeFeedPost{
		TextContent: "This is my first test feed post!",
	}
	userID := 1
	linkSuffix := "Gg4123p893q"
	newPostID, postDbErr := feedrepo.UpsertFeedPostDB(ctx, postReq, userID, response.PostImage{}, linkSuffix, db)
	if postDbErr != nil {
		vlog.Fatalf(ctx, "Error upserting feed post: %s", postDbErr.Error())
	}

	MakeFeedPostComment(db, newPostID, userID, "This is my first test comment on my first test feed post!")
	MakeFeedPostComment(db, newPostID, userID, "This is my second test comment on my first test feed post!")

	postReq2 := request.MakeFeedPost{
		TextContent: "Second test feed post coming through!",
	}
	linkSuffix2 := "Hh4123r893i"
	_, postDbErr2 := feedrepo.UpsertFeedPostDB(ctx, postReq2, userID, response.PostImage{}, linkSuffix2, db)
	if postDbErr2 != nil {
		vlog.Fatalf(ctx, "Error upserting feed post: %s", postDbErr2.Error())
	}
}

func MakeFeedPostComment(db *pgxpool.Pool, postID int, userID int, text string) {
	ctx := context.Background()
	commentReq := request.MakeComment{
		Text:   text,
		PostID: postID,
	}

	_, commentDbErr := feedrepo.MakeCommentDB(context.Background(), commentReq, userID, db)
	if commentDbErr != nil {
		vlog.Fatalf(ctx, "Error inserting comment: %s", commentDbErr.Error())
	}
}

func FillBootstrapContactData(db *pgxpool.Pool) {
	userID := 1
	contactID := 2
	ctx := context.Background()

	contactDBErr := contactrepo.CreateContactDB(ctx, userID, contactID, db)
	if contactDBErr != nil {
		vlog.Fatalf(ctx, "Error creating contact: %s", contactDBErr.Error())
	}

	contactID2 := 3
	contactDBErr2 := contactrepo.CreateContactDB(ctx, userID, contactID2, db)
	if contactDBErr2 != nil {
		vlog.Fatalf(ctx, "Error creating contact: %s", contactDBErr.Error())
	}
}

func FillBoostrapCreateGoatInviteCodeData(db *pgxpool.Pool, inviteCode string) {
	ctx := context.Background()
	createGoatInviteCodeErr := userrepo.CreateGoatInviteCodeDB(ctx, inviteCode, db)
	if createGoatInviteCodeErr != nil {
		vlog.Fatalf(ctx, "Error inserting creator invite code: %s", createGoatInviteCodeErr.Error())
	}
}

func FillBootstrapWalletBalanceData(db *pgxpool.Pool, userID int, balanceDelta int64) {
	ctx := context.Background()
	tx, txErr := db.Begin(ctx)
	if txErr != nil {
		vlog.Fatalf(ctx, "Error creating transaction for FillBootstrapWalletBalanceData: %s", txErr.Error())
	}

	ledgerEntry := wallet.LedgerEntry{
		ProviderUserID: userID,
		Currency:       "usd",
		BalanceDelta:   balanceDelta,
	}

	updateBalanceErr := walletrepo.UpsertBalance(ctx, tx, ledgerEntry)
	if updateBalanceErr != nil {
		vlog.Fatalf(ctx, "Error updating balance : %s", updateBalanceErr.Error())
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		vlog.Fatalf(ctx, "Error committing transaction: %s", commitErr.Error())
	}
}

func FillBootstrapTransactionItems(db *pgxpool.Pool, customerUserID int, providerUserID int) {
	ctx := context.Background()
	tx, txErr := db.Begin(ctx)
	if txErr != nil {
		vlog.Fatalf(ctx, "Error creating transaction for FillBootstrapTransactionItems: %s", txErr.Error())
	}

	totalAmount := int64(1000)
	totalFeesExact := float64(totalAmount) * constants.TOTAL_FEES_RATIO
	stripeFee := int64((float64(totalAmount) * .029) + 30)
	vamaFee := int64(math.Ceil(totalFeesExact)) - int64(stripeFee)

	ledgerEntry1 := wallet.LedgerEntry{
		ProviderUserID:      providerUserID,
		CustomerUserID:      customerUserID,
		Currency:            "usd",
		Amount:              totalAmount,
		StripeTransactionID: "tx_123456789",
		SourceType:          "charge",
		CreatedTS:           time.Now().Unix(),
		StripeFee:           stripeFee,
		VamaFee:             vamaFee,
	}

	insertErr := walletrepo.InsertLedgerTransactionTx(ctx, tx, ledgerEntry1)
	if insertErr != nil {
		vlog.Fatalf(ctx, "Error upserting transaction: %s", insertErr.Error())
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		vlog.Fatalf(ctx, "Error committing transaction: %s", commitErr.Error())
	}
}

func FillBootstrapFollowData(db *pgxpool.Pool) {
	userID := 2
	goatUserID := 1
	ctx := context.Background()
	followErr := followrepo.FollowDB(ctx, userID, goatUserID, db)
	if followErr != nil {
		vlog.Fatalf(ctx, "Error following: %s", followErr.Error())
	}
}

func FillBoostrapUserWithCustomID(runnable utils.Runnable, user response.User) error {
	ctx := context.Background()
	query, args, squirrelErr := squirrel.Insert("core.users").
		Columns(
			"id",
			"first_name",
			"last_name",
			"phone_number",
			"email",
			"username",
			"user_type",
			"stripe_id",
		).
		Values(
			user.ID,
			user.FirstName,
			user.LastName,
			user.Phonenumber,
			user.Email,
			user.Username,
			user.Type,
			user.StripeID,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := runnable.Exec(context.Background(), query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}
	return nil
}

func FillBootstrapVamaUser(runnable utils.Runnable) error {
	vamaUserEmail := "support@vama.com"
	vamaUserStripeID := "Vama0"
	vamaUser := response.User{
		ID:          -1,
		FirstName:   "Vama",
		LastName:    "",
		Phonenumber: "0",
		CountryCode: "US",
		Email:       &vamaUserEmail,
		Username:    "Vama",
		Type:        "VAMA",
		StripeID:    &vamaUserStripeID,
	}
	insertVamaUserErr := FillBoostrapUserWithCustomID(runnable, vamaUser)
	if insertVamaUserErr != nil {
		return insertVamaUserErr
	}

	return nil
}
