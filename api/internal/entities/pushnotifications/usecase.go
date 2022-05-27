package push

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/jackc/pgx/v4"
)

type Usecase interface {
	UpdateUserFcmToken(ctx context.Context, token string) error
	SendTestNotification(ctx context.Context) error

	SetGoatPostNotifications(ctx context.Context, goatID int, enableNotifications bool) error
	AreGoatPostNotificationsEnabled(ctx context.Context, goatUserID int) (bool, error)
	SendGoatPostNotification(ctx context.Context, goatUserID int, postID int) error
	SendRemovedGroupNotification(ctx context.Context, removedUser *response.User, removerUser *response.User, channel *sendbird.GroupChannel) error
	SendPendingBalanceNotifications(ctx context.Context, notifications []wallet.PendingBalanceNotification)

	InitializeSettings(ctx context.Context, userType string, userID int, tx pgx.Tx) error
	UpdateSetting(ctx context.Context, req request.UpdatePushSetting) error
	GetSettings(ctx context.Context) (*response.GetPushSettings, error)
}
