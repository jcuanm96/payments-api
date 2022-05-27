package push

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

type Repository interface {
	baserepo.BaseRepository
	GetUserFcmToken(ctx context.Context, userID int) (string, error)
	UpsertUserFcmToken(ctx context.Context, userID int, token string) error

	EnableGoatPostNotifications(ctx context.Context, runnable utils.Runnable, userID int, goatID int) error
	DisableGoatPostNotifications(ctx context.Context, runnable utils.Runnable, userID int, goatID int) error
	AreGoatPostNotificationsEnabled(ctx context.Context, userID, goatID int) (bool, error)

	GetSettings(ctx context.Context, runnable utils.Runnable, userID int) (*response.GetPushSettings, error)
	UpsertSettings(ctx context.Context, runnable utils.Runnable, userID int, settings UpdateSettings) error
	UpdateSetting(ctx context.Context, runnable utils.Runnable, userID int, req request.UpdatePushSetting) error
}
