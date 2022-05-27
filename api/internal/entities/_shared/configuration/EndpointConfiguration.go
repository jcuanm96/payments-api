package configuration

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/vamawebsocket"
	"github.com/gofiber/fiber/v2"
)

type EndpointConfiguration struct {
	Method           string
	Path             string
	Version          string
	Handler          func(uc interface{}, ctx context.Context, req interface{}) (interface{}, error)
	WebsocketHandler func(uc interface{}, ctx context.Context, c *vamawebsocket.Conn, incomeRequest interface{}) error
	RequestDecoder   func(c *fiber.Ctx) (interface{}, error)
	ResponseEncoder  func(c *fiber.Ctx, response interface{}) error
}

func SuccessResponse(c *fiber.Ctx, response interface{}) error {
	return c.JSON(fiber.Map{"success": true})
}
