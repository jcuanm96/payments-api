package app

import (
	"strings"
	"time"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/controller"
	sconfiguration "github.com/VamaSingapore/vama-api/internal/entities/_shared/configuration"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	authctr "github.com/VamaSingapore/vama-api/internal/entities/auth/controller"
	chatctr "github.com/VamaSingapore/vama-api/internal/entities/chat/controller"
	contactctr "github.com/VamaSingapore/vama-api/internal/entities/contact/controller"
	feedctr "github.com/VamaSingapore/vama-api/internal/entities/feed/controller"
	followctr "github.com/VamaSingapore/vama-api/internal/entities/follow/controller"
	monitoringctr "github.com/VamaSingapore/vama-api/internal/entities/monitoring/controller"
	pushctr "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications/controller"
	searchctr "github.com/VamaSingapore/vama-api/internal/entities/search/controller"
	sharingctr "github.com/VamaSingapore/vama-api/internal/entities/sharing/controller"
	subsriptionctr "github.com/VamaSingapore/vama-api/internal/entities/subscription/controller"
	userctr "github.com/VamaSingapore/vama-api/internal/entities/user/controller"
	walletctr "github.com/VamaSingapore/vama-api/internal/entities/wallet/controller"
	webhooksctr "github.com/VamaSingapore/vama-api/internal/entities/webhooks/controller"
	"github.com/VamaSingapore/vama-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
)

func Setup(app *fiber.App, ctr *controller.Ctr, middlewareCallback fiber.Handler, gcpServiceMiddlewareCallback fiber.Handler) {
	limiter := middleware.Limiter(&middleware.LimiterConfig{Max: 1000, Expiration: time.Second})
	// Global middleware
	app.Use(recover.New())
	if appconfig.Config.Logger {
		app.Use(logger.New())
	}
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	//old endpoints
	authv1 := app.Group("/auth/v1")
	apiv1 := app.Group("/api/v1", middlewareCallback)
	apiv2 := app.Group(constants.API_V2, middlewareCallback)
	websocketGroup := app.Group(constants.WEBSOCKET_V1, middlewareCallback)
	cloudTaskGroup := app.Group(constants.CLOUD_TASK_V1, gcpServiceMiddlewareCallback)
	cloudSchedulerGroup := app.Group(constants.CLOUD_SCHEDULER_V1, gcpServiceMiddlewareCallback)
	monitoring := app.Group("/monitoring")
	webhooks := app.Group("/webhook")
	public := app.Group("/public")

	//declarative endpoint configuration
	{
		usecaseEndpointsMap := make(map[interface{}][]sconfiguration.EndpointConfiguration)

		usecaseEndpointsMap[ctr.Subscription] = subsriptionctr.ComposeEndpoints()
		usecaseEndpointsMap[ctr.User] = userctr.ComposeEndpoints()
		usecaseEndpointsMap[ctr.Wallet] = walletctr.ComposeEndpoints()
		usecaseEndpointsMap[ctr.Feed] = feedctr.ComposeEndpoints()
		usecaseEndpointsMap[ctr.Contact] = contactctr.ComposeEndpoints()
		usecaseEndpointsMap[ctr.Follow] = followctr.ComposeEndpoints()
		usecaseEndpointsMap[ctr.Chat] = chatctr.ComposeEndpoints()
		usecaseEndpointsMap[ctr.Push] = pushctr.ComposeEndpoints()
		usecaseEndpointsMap[ctr.Sharing] = sharingctr.ComposeEndpoints()
		usecaseEndpointsMap[ctr.Search] = searchctr.ComposeEndpoints()

		for uc, endpoints := range usecaseEndpointsMap {
			for _, v := range endpoints {
				switch v.Version {
				case constants.API_V2:
					apiv2.Add(strings.ToUpper(v.Method), v.Path, handlerWrapper(*ctr, v.Handler, v.RequestDecoder, v.ResponseEncoder, uc))
				case constants.PUBLIC_V1:
					public.Add(strings.ToUpper(v.Method), v.Path, handlerWrapper(*ctr, v.Handler, v.RequestDecoder, v.ResponseEncoder, uc))
				case constants.WEBSOCKET_V1:
					websocketGroup.Add(strings.ToUpper(v.Method), v.Path, websocketHandlerWrapper(*ctr, v.WebsocketHandler, v.RequestDecoder, uc))
				case constants.CLOUD_TASK_V1:
					cloudTaskGroup.Add(strings.ToUpper(v.Method), v.Path, handlerWrapper(*ctr, v.Handler, v.RequestDecoder, v.ResponseEncoder, uc))
				case constants.CLOUD_SCHEDULER_V1:
					cloudSchedulerGroup.Add(strings.ToUpper(v.Method), v.Path, handlerWrapper(*ctr, v.Handler, v.RequestDecoder, v.ResponseEncoder, uc))
				default:
					apiv1.Add(strings.ToUpper(v.Method), v.Path, handlerWrapper(*ctr, v.Handler, v.RequestDecoder, v.ResponseEncoder, uc))
				}
			}
		}
	}
	{
		usecaseEndpointsMap := make(map[interface{}][]sconfiguration.EndpointConfiguration)
		usecaseEndpointsMap[ctr.Monitoring] = monitoringctr.ComposeEndpoints()
		for uc, endpoints := range usecaseEndpointsMap {
			for _, v := range endpoints {
				monitoring.Add(strings.ToUpper(v.Method), v.Path, handlerWrapper(*ctr, v.Handler, v.RequestDecoder, v.ResponseEncoder, uc))
			}
		}
	}
	{
		usecaseEndpointsMap := make(map[interface{}][]sconfiguration.EndpointConfiguration)
		usecaseEndpointsMap[ctr.Auth] = authctr.ComposeEndpoints()
		for uc, endpoints := range usecaseEndpointsMap {
			for _, v := range endpoints {
				authv1.Add(strings.ToUpper(v.Method), v.Path, limiter, handlerWrapper(*ctr, v.Handler, v.RequestDecoder, v.ResponseEncoder, uc))
			}
		}
	}
	{
		usecaseEndpointsMap := make(map[interface{}][]sconfiguration.EndpointConfiguration)
		usecaseEndpointsMap[ctr.Webhooks] = webhooksctr.ComposeEndpoints()
		for uc, endpoints := range usecaseEndpointsMap {
			for _, v := range endpoints {
				webhooks.Add(strings.ToUpper(v.Method), v.Path, handlerWrapper(*ctr, v.Handler, v.RequestDecoder, v.ResponseEncoder, uc))
			}
		}
	}
	{
		apiv1.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("v1")
		})
	}
}
