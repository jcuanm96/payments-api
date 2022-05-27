package main

import (
	"context"
	"fmt"
	"os"

	repo "github.com/VamaSingapore/vama-api/cmd/scripts/upgradeUserToGoat/repository"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sharedrepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	userrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
)

func main() {
	userUpdateEmail := ""
	userUpdatePhone := ""
	userUpdateFields := response.User{
		Phonenumber: userUpdatePhone,
		Email:       &userUpdateEmail,
		Type:        "GOAT",
	}

	appconfig.Init("configs/config.yaml")
	// tmp code injection
	connectionString := appconfig.Config.Postgre.ConnectionString
	pgxConfig, configErr := pgxpool.ParseConfig(connectionString)
	if configErr != nil {
		logrus.Fatalf("Error parsing db connection string %s: %s", connectionString, configErr.Error())
	}
	ctx := context.Background()
	pool, dbErr := pgxpool.ConnectConfig(ctx, pgxConfig)
	if dbErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbErr)
		logrus.Fatalf(dbErr.Error())
	}
	defer pool.Close()

	tx, txErr := pool.Begin(ctx)
	if txErr != nil {
		logrus.Fatalf("Error beginning tx: %s", txErr.Error())
	}

	commit := false
	defer sharedrepo.FinishTx(ctx, tx, &commit)

	// Client initialization
	stripeClient := &client.API{}
	stripeBackendConfig := &stripe.BackendConfig{
		MaxNetworkRetries: stripe.Int64(3), // Automatically create an idempotency key for all POST requests to Stripe and retry on API failures
	}
	stripeClient.Init(appconfig.Config.Stripe.Key, &stripe.Backends{
		API: stripe.GetBackendWithConfig(stripe.APIBackend, stripeBackendConfig),
	})

	var rdb vredis.Client
	userrepoClient := userrepo.New(pool, rdb)

	currUser, getUserByPhoneErr := userrepoClient.GetUserByPhone(ctx, tx, userUpdateFields.Phonenumber)
	if getUserByPhoneErr != nil {
		logrus.Fatalf("Error accessing db for phone number %s. Err: %s\n", userUpdateFields.Phonenumber, getUserByPhoneErr.Error())
	} else if currUser == nil {
		logrus.Fatalf("The user doesn't exist yet. Please have them make an account before upgrading")
	}

	updateErr := repo.UpdateUser(ctx, tx, userUpdateFields.Phonenumber, *userUpdateFields.Email)
	if updateErr != nil {
		logrus.Fatalf("Error accessing db for phone number %s. Err: %s\n", userUpdateFields.Phonenumber, updateErr.Error())
	}

	updateStripeCustomerParams := &stripe.CustomerParams{
		Email: stripe.String(*userUpdateFields.Email),
		Phone: stripe.String(userUpdateFields.Phonenumber),
	}

	_, updateStripeCustomerErr := stripeClient.Customers.Update(
		*currUser.StripeID,
		updateStripeCustomerParams,
	)

	if updateStripeCustomerErr != nil {
		logrus.Fatalf("Could not update Stripe user for phone number %s. Error: %s\n", userUpdateFields.Phonenumber, updateStripeCustomerErr.Error())
	}

	defaultCurrency := constants.DEFAULT_CURRENCY
	defaultGoatChatPrice := int64(constants.MIN_GOAT_CHAT_PRICE_USD_IN_SMALLEST_DENOM)

	// Initialize a GOATs chat price
	upsertGoatChatsPriceErr := userrepoClient.UpsertGoatChatsPrice(ctx, tx, defaultGoatChatPrice, defaultCurrency, currUser.ID)
	if upsertGoatChatsPriceErr != nil {
		logrus.Fatalf("Error initializing creator chat price for phone %s. Error: %s\n", userUpdateFields.Phonenumber, upsertGoatChatsPriceErr.Error())
	}

	sendbirdClient := sendbird.NewClient()
	sendBirdMetadataReq := request.SendBirdUserMetadata{
		Type:     "GOAT",
		Username: request.Username{Username: currUser.Username},
	}
	sendBirdMetadataErr := sendbirdClient.UpsertUserMetadata(currUser.ID, sendBirdMetadataReq)
	if sendBirdMetadataErr != nil {
		logrus.Fatalf("Error associating username with sendbird user %d. Error: %s\n", currUser.ID, sendBirdMetadataErr.Error())
	}

	commit = true
}
