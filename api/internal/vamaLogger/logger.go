package vlog

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/logging"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

var VamaLoggerClient *logging.Client

func Init(ctx context.Context) error {
	cloudLoggingClient, newClientErr := logging.NewClient(ctx, appconfig.Config.Gcloud.Project)
	if newClientErr != nil {
		return newClientErr
	}

	VamaLoggerClient = cloudLoggingClient

	return nil
}

func log(ctx context.Context,
	errMsg string,
	severity logging.Severity,
	localLogFunc func(args ...interface{})) {
	if VamaLoggerClient == nil {
		telegram.TelegramClient.SendMessage("VamaLoggerClient is nil")
		return
	}

	var userID string
	if ctx.Value("CURRENT_USER_UUID") != nil {
		userID = ctx.Value("CURRENT_USER_UUID").(string)
	}

	if appconfig.Config.VamaLogger.Type != constants.GCP_LOGGER_TYPE {
		localLogFunc(errMsg)
		return
	}

	var logID string
	if ctx.Value("REQUEST_ID") != nil {
		logID = ctx.Value("REQUEST_ID").(string)
	} else {
		uuid, uuidErr := uuid.NewV4()
		if uuidErr != nil {
			logrus.Error("Could not get uuid for logger. Err: %v", uuidErr)
			return
		}
		logID = uuid.String()
	}

	logger := VamaLoggerClient.Logger(logID)
	logger.Log(logging.Entry{
		Severity: severity,
		Payload:  errMsg,
		Labels: map[string]string{
			"userID": userID,
		},
	})
}

func Error(ctx context.Context, errMsg string) {
	log(ctx, errMsg, logging.Warning, logrus.Error)
}

func Errorf(ctx context.Context, errMsgf string, args ...interface{}) {
	errMsg := fmt.Sprintf(errMsgf, args...)
	Error(ctx, errMsg)
}

func Fatal(ctx context.Context, errMsg string) {
	log(ctx, errMsg, logging.Critical, logrus.Fatal)
	os.Exit(1)
}

func Fatalf(ctx context.Context, errMsgf string, args ...interface{}) {
	errMsg := fmt.Sprintf(errMsgf, args...)
	Fatal(ctx, errMsg)
}

func Info(ctx context.Context, msg string) {
	log(ctx, msg, logging.Info, logrus.Info)
}

func Infof(ctx context.Context, msgf string, args ...interface{}) {
	msg := fmt.Sprintf(msgf, args...)
	Info(ctx, msg)
}

func Close() {
	if VamaLoggerClient != nil {
		VamaLoggerClient.Close()
	}
}
