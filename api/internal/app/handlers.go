package app

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/controller"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/internal/vamawebsocket"
	"github.com/gofiber/fiber/v2"
)

var silentFlag = flag.Bool("silent", false, "disable request/response logging")

func handlerWrapper(
	ctr controller.Ctr,
	handler func(uc interface{}, ctx context.Context, req interface{}) (interface{}, error),
	requestDecoder func(c *fiber.Ctx) (interface{}, error),
	responseEncoder func(c *fiber.Ctx, response interface{}) error,
	uc interface{},
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var reqDecoded interface{}
		var err error
		ctx, prepareCtxErr := ctr.PrepareContext(c)
		if prepareCtxErr != nil {
			errMsg := fmt.Sprintf("error getting context. Err: %v", prepareCtxErr)
			telegram.TelegramClient.SendMessage(errMsg)
			return ctr.WrapError(errors.New(errMsg), c)
		} else if ctx == nil {
			errMsg := "error context came back as nil"
			telegram.TelegramClient.SendMessage(errMsg)
			return ctr.WrapError(errors.New(errMsg), c)
		}

		if requestDecoder != nil {
			reqDecoded, err = requestDecoder(c)
			if err != nil {
				return ctr.WrapError(err, c)
			}
			if reqDecoded == nil {
				return ctr.WrapError(errors.New("error parsing request"), c)
			}
		}

		if !*silentFlag {
			reqBytes, marshalErr := json.Marshal(reqDecoded)
			if marshalErr != nil {
				vlog.Errorf(*ctx, "Error marshaling request to log: %s", marshalErr)
			}
			reqString := string(reqBytes)
			if reqString == "null" {
				reqString = "{}"
			}
			vlog.Infof(*ctx, "Request: %s", reqString)
		}
		res, err := handler(uc, *ctx, reqDecoded)

		if err != nil {
			vlog.Errorf(*ctx, "Response error: %v", err)
			return ctr.WrapError(err, c)
		}

		if res == nil {
			res = struct{}{}
		}

		if responseEncoder != nil {
			vlog.EndpointInfo(c)
			return responseEncoder(c, res)
		}

		jsonErr := c.JSON(res)
		if !*silentFlag {
			vlog.EndpointInfo(c)
		}
		return jsonErr
	}
}

func websocketHandlerWrapper(
	ctr controller.Ctr,
	handler func(uc interface{}, ctx context.Context, c *vamawebsocket.Conn, incomeRequest interface{}) error,
	requestDecoder func(c *fiber.Ctx) (interface{}, error),
	uc interface{},
) func(*fiber.Ctx) error {
	return vamawebsocket.New(func(conn *vamawebsocket.Conn) {
		c := conn.Ctx()
		var reqDecoded interface{}
		var decodeReqErr error

		ctx, prepareCtxErr := ctr.PrepareContext(c)
		if prepareCtxErr != nil {
			errMsg := fmt.Sprintf("error getting context. Err: %v", prepareCtxErr)
			telegram.TelegramClient.SendMessage(errMsg)
			return
		} else if ctx == nil {
			errMsg := "error context came back as nil"
			telegram.TelegramClient.SendMessage(errMsg)
			return
		}

		if requestDecoder != nil {
			reqDecoded, decodeReqErr = requestDecoder(c)
			if decodeReqErr != nil {
				return
			}
			if reqDecoded == nil {
				return
			}
		}

		if !*silentFlag {
			reqBytes, marshalErr := json.Marshal(reqDecoded)
			if marshalErr != nil {
				vlog.Errorf(*ctx, "Error marshaling request to log: %s", marshalErr)
			}
			reqString := string(reqBytes)
			if reqString == "null" {
				reqString = "{}"
			}
			vlog.Infof(*ctx, "Websocket request: %s", reqString)
		}

		handlerErr := handler(uc, *ctx, conn, reqDecoded)
		if handlerErr != nil {
			return
		}
	})
}
