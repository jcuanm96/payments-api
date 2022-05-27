package main

import (
	"context"
	"fmt"
	"os"

	repo "github.com/VamaSingapore/vama-api/cmd/scripts/createVamaUser/repository"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sharedrepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func main() {
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

	asUserID := -1
	email := "support@vama.com"
	stripeID := "Vama0"
	vamaUser := &response.User{
		ID:          asUserID,
		FirstName:   "Vama",
		Username:    "Vama",
		Phonenumber: "0",
		Type:        constants.VAMA_USER_TYPE,
		Email:       &email,
		StripeID:    &stripeID,
	}

	upsertErr := repo.UpsertVamaUser(tx, vamaUser)
	if upsertErr != nil {
		logrus.Fatalf("Error upserting vama user %v: %s", *vamaUser, upsertErr.Error())
	}

	sendbirdClient := sendbird.NewClient()

	createUserParams := &sendbird.CreateUserParams{
		UserID:   asUserID,
		Nickname: "Vama",
	}
	_, createSendbirdUserErr := sendbirdClient.CreateUser(createUserParams)
	if createSendbirdUserErr != nil {
		logrus.Fatalf("Error creating sendbird user %d: %s", asUserID, createSendbirdUserErr.Error())
	}

	commit = true
}
