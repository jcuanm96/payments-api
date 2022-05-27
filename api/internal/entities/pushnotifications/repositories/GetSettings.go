package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (s *repository) GetSettings(ctx context.Context, runnable utils.Runnable, userID int) (*response.GetPushSettings, error) {
	query, args, squirrelErr := squirrel.Select(
		"pending_balance",
	).
		From("push.settings").
		Where("user_id = ?", userID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if squirrelErr != nil {
		return nil, squirrelErr
	}

	row := runnable.QueryRow(ctx, query, args...)

	pendingBalanceSetting := response.PushSetting{
		ID:       push.SettingConfigs.PendingBalance.ID,
		Title:    push.SettingConfigs.PendingBalance.Title,
		Category: push.SettingConfigs.PendingBalance.Category,
	}
	scanErr := row.Scan(
		&pendingBalanceSetting.Setting,
	)

	if scanErr != nil {
		return nil, scanErr
	}

	settings := []response.PushSetting{
		pendingBalanceSetting,
	}
	res := &response.GetPushSettings{
		Settings: settings,
	}
	return res, nil

}
