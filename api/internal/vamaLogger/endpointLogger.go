package vlog

import (
	"fmt"

	"cloud.google.com/go/logging"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func EndpointError(fiberCtx *fiber.Ctx, status int, detail string) {
	if VamaLoggerClient == nil {
		telegram.TelegramClient.SendMessage("VamaLoggerClient is nil")
		return
	}

	var userID string
	if fiberCtx.Locals("uid") != nil {
		userID = fiberCtx.Locals("uid").(string)
	}

	if appconfig.Config.VamaLogger.Type != constants.GCP_LOGGER_TYPE {
		logrus.Error(detail)
		return
	}

	var logID string
	requestID, ok := fiberCtx.Locals("REQUEST_ID").(string)
	if ok && requestID != "" {
		logID = requestID
	} else {
		uuid, uuidErr := uuid.NewV4()
		if uuidErr != nil {
			logrus.Error("Could not get uuid for logger. Err: %v", uuidErr)
			return
		}
		logID = uuid.String()
	}

	requestParams := ""
	queryParams := string(fiberCtx.Context().QueryArgs().QueryString())
	if queryParams != "" {
		requestParams += "query: " + queryParams + "\n"
	}
	bodyParams := string(fiberCtx.Body())
	if bodyParams != "" {
		requestParams += "body: " + bodyParams
	}

	logger := VamaLoggerClient.Logger(logID)
	path := string(fiberCtx.Context().Path())
	method := string(fiberCtx.Context().Method())
	payload := fmt.Sprintf("%s | %s | %d", path, method, status)
	logger.Log(logging.Entry{
		Severity: logging.Error,
		Payload:  payload,
		Labels: map[string]string{
			"requestParams": requestParams,
			"requestMethod": method,
			"errMsg":        detail,
			"statusCode":    fmt.Sprint(status),
			"responseBody":  string(fiberCtx.Response().Body()),
			"userID":        userID,
		},
	})
}

func EndpointInfo(fiberCtx *fiber.Ctx) {
	if VamaLoggerClient == nil {
		telegram.TelegramClient.SendMessage("VamaLoggerClient is nil")
		return
	}

	path := fiberCtx.Context().Path()
	method := string(fiberCtx.Context().Method())
	payload := fmt.Sprintf("%s | %s", method, path)

	if appconfig.Config.VamaLogger.Type != constants.GCP_LOGGER_TYPE {
		logrus.Info(string(payload))
		return
	}

	var logID string
	requestID, ok := fiberCtx.Locals("REQUEST_ID").(string)
	if ok && requestID != "" {
		logID = requestID
	} else {
		uuid, uuidErr := uuid.NewV4()
		if uuidErr != nil {
			logrus.Error("Could not get uuid for logger. Err: %v", uuidErr)
			return
		}
		logID = uuid.String()
	}

	var userID string
	if fiberCtx.Locals("uid") != nil {
		userID = fiberCtx.Locals("uid").(string)
	}

	logger := VamaLoggerClient.Logger(logID)
	logger.Log(logging.Entry{
		Severity: logging.Info,
		Payload:  payload,
		Labels: map[string]string{
			"requestParams": string(fiberCtx.Request().Body()),
			"requestMethod": string(fiberCtx.Context().Method()),
			"responseBody":  string(fiberCtx.Response().Body()),
			"userID":        userID,
		},
	})
}
