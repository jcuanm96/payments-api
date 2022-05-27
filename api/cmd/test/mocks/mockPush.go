package mocks

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/jackc/pgx/v4"
)

type MockPushUsecase struct{}

func NewMockPush() push.Usecase {
	return &MockPushUsecase{}
}

func (muc *MockPushUsecase) UpdateUserFcmToken(ctx context.Context, token string) error {
	return nil
}
func (muc *MockPushUsecase) SendTestNotification(ctx context.Context) error {
	return nil
}

func (muc *MockPushUsecase) SetGoatPostNotifications(ctx context.Context, goatID int, enableNotifications bool) error {
	return nil
}
func (muc *MockPushUsecase) AreGoatPostNotificationsEnabled(ctx context.Context, goatUserID int) (bool, error) {
	return false, nil
}

func (muc *MockPushUsecase) SendGoatPostNotification(ctx context.Context, goatUserID int, postID int) error {
	return nil
}

func (muc *MockPushUsecase) SendRemovedGroupNotification(ctx context.Context, removedUser *response.User, removerUser *response.User, channel *sendbird.GroupChannel) error {
	return nil
}
func (muc *MockPushUsecase) SendPendingBalanceNotifications(ctx context.Context, notifications []wallet.PendingBalanceNotification) {
}

func (muc *MockPushUsecase) UpdateSetting(ctx context.Context, req request.UpdatePushSetting) error {
	return nil
}

func (muc *MockPushUsecase) InitializeSettings(ctx context.Context, userType string, userID int, tx pgx.Tx) error {
	return nil
}

func (muc *MockPushUsecase) GetSettings(ctx context.Context) (*response.GetPushSettings, error) {
	return nil, nil
}
