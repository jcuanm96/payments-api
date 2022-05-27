package test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/VamaSingapore/vama-api/internal/entities/auth"
	authrepo "github.com/VamaSingapore/vama-api/internal/entities/auth/repositories"
	authusecase "github.com/VamaSingapore/vama-api/internal/entities/auth/usecase"
	"github.com/VamaSingapore/vama-api/internal/entities/chat"
	"github.com/VamaSingapore/vama-api/internal/entities/follow"
	"github.com/VamaSingapore/vama-api/internal/entities/search"
	"github.com/VamaSingapore/vama-api/internal/entities/sharing"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	"github.com/VamaSingapore/vama-api/internal/messaging"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	"github.com/VamaSingapore/vama-api/internal/vamabot"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	chatrepo "github.com/VamaSingapore/vama-api/internal/entities/chat/repositories"
	chatusecase "github.com/VamaSingapore/vama-api/internal/entities/chat/usecase"

	contactrepo "github.com/VamaSingapore/vama-api/internal/entities/contact/repositories"
	contactusecase "github.com/VamaSingapore/vama-api/internal/entities/contact/usecase"

	feedrepo "github.com/VamaSingapore/vama-api/internal/entities/feed/repositories"
	feedusecase "github.com/VamaSingapore/vama-api/internal/entities/feed/usecase"

	followrepo "github.com/VamaSingapore/vama-api/internal/entities/follow/repositories"
	followusecase "github.com/VamaSingapore/vama-api/internal/entities/follow/usecase"

	userrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	userusecase "github.com/VamaSingapore/vama-api/internal/entities/user/usecase"

	monitoringrepo "github.com/VamaSingapore/vama-api/internal/entities/monitoring/repositories"
	monitoringusecase "github.com/VamaSingapore/vama-api/internal/entities/monitoring/usecase"

	searchrepo "github.com/VamaSingapore/vama-api/internal/entities/search/repositories"
	searchusecase "github.com/VamaSingapore/vama-api/internal/entities/search/usecase"

	webhooksrepo "github.com/VamaSingapore/vama-api/internal/entities/webhooks/repositories"
	webhooksusecase "github.com/VamaSingapore/vama-api/internal/entities/webhooks/usecase"

	subscriptionrepo "github.com/VamaSingapore/vama-api/internal/entities/subscription/repositories"
	subscriptionusecase "github.com/VamaSingapore/vama-api/internal/entities/subscription/usecase"

	walletrepo "github.com/VamaSingapore/vama-api/internal/entities/wallet/repositories"
	walletusecase "github.com/VamaSingapore/vama-api/internal/entities/wallet/usecase"

	"github.com/VamaSingapore/vama-api/internal/app"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/controller"
	"github.com/VamaSingapore/vama-api/internal/logerr"
	"github.com/VamaSingapore/vama-api/internal/token"
	"github.com/VamaSingapore/vama-api/internal/upload"
	"github.com/gofiber/fiber/v2"
	"github.com/kevinburke/twilio-go"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"

	mock "github.com/VamaSingapore/vama-api/cmd/test/mocks"
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

type TestApp struct {
	App   *fiber.App
	Ctr   *controller.Ctr
	Db    *pgxpool.Pool
	Redis vredis.Client
}

func StartTestServer() TestApp {
	// Set path as project root
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	flag.Set("silent", "true")
	flag.Parse()
	appconfig.Init("configs/test-config.yaml")
	ctx := context.Background()
	vlogInitErr := vlog.Init(ctx)
	if vlogInitErr != nil {
		logrus.Fatalf("Error initializing vama logger: %v", vlogInitErr)
	}
	defer vlog.Close()

	telegram.Init()

	var lr *logerr.Logerr

	connectionString := appconfig.Config.Postgre.ConnectionString
	if connectionString != "host=localhost user=postgres password=pass database=postgres-test sslmode=disable" {
		vlog.Fatalf(ctx, "unexpected database string for test DB: %s", connectionString)
	}
	pgxConfig, configErr := pgxpool.ParseConfig(connectionString)
	if configErr != nil {
		vlog.Fatalf(ctx, "Error parsing db connection string %s: %s", connectionString, configErr.Error())
	}
	pool, dbErr := pgxpool.ConnectConfig(context.Background(), pgxConfig)
	if dbErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbErr)
		vlog.Fatalf(ctx, dbErr.Error())
	}

	rdb := mock.NewRedisClient()

	// Setup application controller with its dependencies.
	var sharinguc sharing.Usecase

	msg := messaging.New()
	sendbirdClient := mock.NewSendBirdMockClient()
	vamaBot := vamabot.NewClient(1, sendbirdClient)
	cloudTasksClient := mock.NewMockCloudTasks()

	var followuc follow.Usecase
	var walletuc wallet.Usecase

	// Configure Stripe client.
	stripeClient := &client.API{}
	stripeBackendConfig := &stripe.BackendConfig{
		MaxNetworkRetries: stripe.Int64(3), // Automatically create an idempotency key for all POST requests to Stripe and retry on API failures
	}
	stripeClient.Init(appconfig.Config.Stripe.Key, &stripe.Backends{
		API: stripe.GetBackendWithConfig(stripe.APIBackend, stripeBackendConfig),
	})

	var twilioClient *twilio.VerifyPhoneNumberService
	var gcsClient *upload.Client
	var authuc auth.Usecase
	var searchuc search.Usecase
	var chatuc chat.Usecase

	// mock usecases
	mockPushuc := mock.NewMockPush()
	realUserUc := userusecase.New(
		userrepo.New(pool, rdb),
		&mockPushuc,
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
	// We pass the real useruc to use the production methods that don't rely on third party services
	// without copy-pasting that code in.
	mockUseruc := mock.NewMockUser(userrepo.New(pool, rdb), realUserUc)

	// real usecases
	tokenSvc := token.NewService(
		pool,
		rdb,
		appconfig.Config.Auth.AccessTokenKey,
	)
	subscriptionuc := subscriptionusecase.New(
		subscriptionrepo.New(pool, rdb),
		realUserUc,
		&chatuc,
		stripeClient,
		sendbirdClient,
	)
	authuc = authusecase.New(
		authrepo.New(pool, rdb),
		realUserUc,
		subscriptionuc,
		tokenSvc,
		msg,
		sendbirdClient,
		vamaBot,
		twilioClient,
		stripeClient,
		cloudTasksClient,
		&walletuc,
		mockPushuc,
	)
	walletuc = walletusecase.New(
		walletrepo.New(pool),
		mockUseruc,
		stripeClient,
		authuc,
		mockPushuc,
	)
	searchuc = searchusecase.New(searchrepo.New(pool), sendbirdClient, mockUseruc)
	contactuc := contactusecase.New(contactrepo.New(pool, rdb), mockUseruc, searchuc)
	monitoringuc := monitoringusecase.New(monitoringrepo.New(pool, rdb))
	feeduc := feedusecase.New(
		feedrepo.New(pool),
		mockUseruc,
		&mockPushuc,
		&chatuc,
		subscriptionuc,
		gcsClient,
	)
	chatuc = chatusecase.New(
		chatrepo.New(pool, rdb),
		feeduc,
		mockUseruc,
		walletuc,
		subscriptionuc,
		msg,
		sendbirdClient,
		stripeClient,
		cloudTasksClient,
		gcsClient,
		mockPushuc,
	)

	webhooksuc := webhooksusecase.New(
		webhooksrepo.New(pool),
		stripeClient,
		sendbirdClient,
		walletuc,
		subscriptionuc,
		chatuc,
	)
	followuc = followusecase.New(followrepo.New(pool), mockUseruc)

	ctr := controller.New(
		walletuc,
		monitoringuc,
		contactuc,
		authuc,
		mockUseruc,
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
		mockPushuc,
		followuc,
		*stripeClient,
		searchuc,
	)

	// Configure and run fiber.
	a := fiber.New()
	p := func(c *fiber.Ctx) error {
		c.Locals("uid", c.Get("id"))
		return c.Next()
	}
	cloudTaskMiddlewareCallback := func(c *fiber.Ctx) error {
		return c.Next()
	}

	app.Setup(a, ctr, p, cloudTaskMiddlewareCallback)
	return TestApp{App: a, Ctr: ctr, Db: pool, Redis: rdb}
}
