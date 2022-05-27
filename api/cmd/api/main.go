package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	cloudTasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/auth"
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	contactrepo "github.com/VamaSingapore/vama-api/internal/entities/contact/repositories"
	contactusecase "github.com/VamaSingapore/vama-api/internal/entities/contact/usecase"
	"github.com/VamaSingapore/vama-api/internal/entities/follow"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/entities/search"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/internal/middleware"
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"

	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/internal/vamabot"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/storage"
	"github.com/VamaSingapore/vama-api/internal/app"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/controller"
	"github.com/VamaSingapore/vama-api/internal/logerr"
	"github.com/VamaSingapore/vama-api/internal/messaging"
	"github.com/VamaSingapore/vama-api/internal/token"
	"github.com/VamaSingapore/vama-api/internal/upload"

	"github.com/gofiber/fiber/v2"
	"github.com/kevinburke/twilio-go"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"

	followrepo "github.com/VamaSingapore/vama-api/internal/entities/follow/repositories"
	followusecase "github.com/VamaSingapore/vama-api/internal/entities/follow/usecase"

	monitoringrepo "github.com/VamaSingapore/vama-api/internal/entities/monitoring/repositories"
	monitoringusecase "github.com/VamaSingapore/vama-api/internal/entities/monitoring/usecase"

	userrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	userusecase "github.com/VamaSingapore/vama-api/internal/entities/user/usecase"

	authrepo "github.com/VamaSingapore/vama-api/internal/entities/auth/repositories"
	authusecase "github.com/VamaSingapore/vama-api/internal/entities/auth/usecase"

	chatrepo "github.com/VamaSingapore/vama-api/internal/entities/chat/repositories"
	chatusecase "github.com/VamaSingapore/vama-api/internal/entities/chat/usecase"

	subscriptionrepo "github.com/VamaSingapore/vama-api/internal/entities/subscription/repositories"
	subscriptionusecase "github.com/VamaSingapore/vama-api/internal/entities/subscription/usecase"

	walletrepo "github.com/VamaSingapore/vama-api/internal/entities/wallet/repositories"
	walletusecase "github.com/VamaSingapore/vama-api/internal/entities/wallet/usecase"

	feedrepo "github.com/VamaSingapore/vama-api/internal/entities/feed/repositories"
	feedusecase "github.com/VamaSingapore/vama-api/internal/entities/feed/usecase"

	pushrepo "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications/repositories"
	pushusecase "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications/usecase"

	sharingrepo "github.com/VamaSingapore/vama-api/internal/entities/sharing/repositories"
	sharingusecase "github.com/VamaSingapore/vama-api/internal/entities/sharing/usecase"

	searchrepo "github.com/VamaSingapore/vama-api/internal/entities/search/repositories"
	searchusecase "github.com/VamaSingapore/vama-api/internal/entities/search/usecase"

	webhooksrepo "github.com/VamaSingapore/vama-api/internal/entities/webhooks/repositories"
	webhooksusecase "github.com/VamaSingapore/vama-api/internal/entities/webhooks/usecase"
)

func main() {
	flag.Parse()
	// Initialize configuration service.
	appconfig.Init("configs/config.yaml")
	ctx := context.Background()

	vlogInitErr := vlog.Init(ctx)
	if vlogInitErr != nil {
		logrus.Fatalf("Error initializing vama logger: %v", vlogInitErr)
	}
	defer vlog.Close()

	telegram.Init()

	// Configure Twilio client.
	twilioClient := twilio.NewClient(appconfig.Config.Twilio.SID, appconfig.Config.Twilio.Token, nil).Verify.Verifications

	// Configure Stripe client.
	stripeClient := &client.API{}
	stripeBackendConfig := &stripe.BackendConfig{
		MaxNetworkRetries: stripe.Int64(3), // Automatically create an idempotency key for all POST requests to Stripe and retry on API failures
	}
	stripeClient.Init(appconfig.Config.Stripe.Key, &stripe.Backends{
		API: stripe.GetBackendWithConfig(stripe.APIBackend, stripeBackendConfig),
	})

	rdb, newRedisClientErr := vredis.New(ctx, appconfig.Config.Redis.Endpoint, appconfig.Config.Redis.Port, appconfig.Config.Redis.Password)
	if newRedisClientErr != nil {
		vlog.Fatalf(ctx, "Error instantiating redis client: %v", newRedisClientErr)
	}

	// Configure Cloud Storage and upload service.
	storageClient, errorReportingErr := storage.NewClient(ctx)
	if errorReportingErr != nil {
		log.Fatalf("failed to create storage service: %v", errorReportingErr)
	}

	projectID := appconfig.Config.Gcloud.Project

	profileAvatarBucket := storageClient.Bucket(appconfig.Config.Storage.ProfileAvatarBucketName)
	feedPostBucket := storageClient.Bucket(appconfig.Config.Storage.FeedPostBucketName)
	themeBucket := storageClient.Bucket(appconfig.Config.Storage.ThemesBucketName)
	chatMediaBucket := storageClient.Bucket(appconfig.Config.Storage.ChatMediaBucketName)
	gcsClient := upload.New(
		profileAvatarBucket,
		feedPostBucket,
		themeBucket,
		chatMediaBucket,

		appconfig.Config.Storage.ProfileAvatarBucketName,
		appconfig.Config.Storage.FeedPostBucketName,
		appconfig.Config.Storage.ThemesBucketName,
		appconfig.Config.Storage.ChatMediaBucketName,
	)

	errorClient, errorReportingErr := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: projectID + "-api",
		OnError: func(err error) {
			log.Fatal(err)
		},
	})
	if errorReportingErr != nil {
		log.Fatalf("failed to create error reporting client: %v", errorReportingErr)
	}
	defer errorClient.Close()

	lr := logerr.New(errorClient)

	// tmp code injection
	connectionString := appconfig.Config.Postgre.ConnectionString
	pgxConfig, configErr := pgxpool.ParseConfig(connectionString)
	if configErr != nil {
		vlog.Fatalf(ctx, "Error parsing db connection string %s: %v", connectionString, configErr)
	}
	pool, dbErr := pgxpool.ConnectConfig(ctx, pgxConfig)
	if dbErr != nil {
		vlog.Fatalf(ctx, "Unable to connect to database: %v", dbErr)
	}
	defer pool.Close()

	// Configure messaging service.
	tokenSvc := token.NewService(pool, rdb, appconfig.Config.Auth.AccessTokenKey)

	msg := messaging.New()
	sendbirdClient := sendbird.NewClient()
	vamaBot := vamabot.NewClient(constants.VAMA_USER_ID, sendbirdClient)
	cloudTasksClient, cloudTasksClientErr := cloudTasks.NewClient(ctx)
	if cloudTasksClientErr != nil {
		vlog.Fatalf(ctx, "Error initializing cloudTasksClient: %v", cloudTasksClientErr)
	}

	var pushuc push.Usecase
	var followuc follow.Usecase
	var chatuc chat.Usecase
	var authuc auth.Usecase
	var walletuc wallet.Usecase
	var searchuc search.Usecase

	useruc := userusecase.New(
		userrepo.New(pool, rdb),
		&pushuc,
		&followuc,
		sendbirdClient,
		msg,
		gcsClient,
		twilioClient,
		appconfig.Config.Twilio.Verify,
		stripeClient,
		&authuc,
		&searchuc,
	)
	pushuc = pushusecase.New(ctx, pushrepo.New(pool), useruc)
	followuc = followusecase.New(followrepo.New(pool), useruc)
	subscriptionuc := subscriptionusecase.New(
		subscriptionrepo.New(pool, rdb),
		useruc,
		&chatuc,
		stripeClient,
		sendbirdClient,
	)
	authuc = authusecase.New(
		authrepo.New(pool, rdb),
		useruc,
		subscriptionuc,
		tokenSvc,
		msg,
		sendbirdClient,
		vamaBot,
		twilioClient,
		stripeClient,
		cloudTasksClient,
		&walletuc,
		pushuc,
	)
	searchuc = searchusecase.New(searchrepo.New(pool), sendbirdClient, useruc)
	contactuc := contactusecase.New(contactrepo.New(pool, rdb), useruc, searchuc)
	monitoringuc := monitoringusecase.New(monitoringrepo.New(pool, rdb))
	walletuc = walletusecase.New(
		walletrepo.New(pool),
		useruc,
		stripeClient,
		authuc,
		pushuc,
	)
	feeduc := feedusecase.New(
		feedrepo.New(pool),
		useruc,
		&pushuc,
		&chatuc,
		subscriptionuc,
		gcsClient,
	)

	chatuc = chatusecase.New(
		chatrepo.New(pool, rdb),
		feeduc,
		useruc,
		walletuc,
		subscriptionuc,
		msg,
		sendbirdClient,
		stripeClient,
		cloudTasksClient,
		gcsClient,
		pushuc,
	)
	sharinguc := sharingusecase.New(sharingrepo.New(pool), useruc, sendbirdClient, gcsClient)
	webhooksuc := webhooksusecase.New(webhooksrepo.New(pool), stripeClient, sendbirdClient, walletuc, subscriptionuc, chatuc)
	// Setup application controller with its dependencies.
	ctr := controller.New(
		walletuc,
		monitoringuc,
		contactuc,
		authuc,
		useruc,
		chatuc,
		subscriptionuc,
		tokenSvc,
		gcsClient,
		msg,
		sendbirdClient,
		lr,
		feeduc,
		sharinguc,
		webhooksuc,
		pushuc,
		followuc,
		*stripeClient,
		searchuc,
	)

	// Configure and run fiber.
	a := fiber.New()
	middlewareCallback := middleware.MustParseClaims(ctr.TokenSvc())
	gcpServiceMiddlewareCallback := middleware.ParseGCPServiceAuthorization(ctr.TokenSvc())
	app.Setup(a, ctr, middlewareCallback, gcpServiceMiddlewareCallback)
	addr := fmt.Sprintf(":%s", appconfig.Config.Port)
	log.Fatal(a.Listen(addr))
}
