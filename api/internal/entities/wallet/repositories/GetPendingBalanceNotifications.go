package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/wallet"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetPendingBalanceNotifications(ctx context.Context, runnable utils.Runnable) ([]wallet.PendingBalanceNotification, error) {
	query, args, squirrelErr := squirrel.Select(
		"balance.available_balance",
		"token.fcm_token",
	).
		From("wallet.balances balance").
		Where("balance.available_balance > 0").
		Where("bank_info.id IS NULL").
		Join("core.tokens token ON balance.provider_user_id = token.user_id").
		LeftJoin("wallet.bank_info bank_info ON balance.provider_user_id = bank_info.user_id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	rows, queryErr := runnable.Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}

	defer rows.Close()

	notifications := []wallet.PendingBalanceNotification{}
	for rows.Next() {
		notification := wallet.PendingBalanceNotification{}
		scanErr := rows.Scan(
			&notification.AvailableBalance,
			&notification.FcmToken,
		)

		if scanErr != nil {
			return nil, scanErr
		}

		notifications = append(notifications, notification)
	}
	return notifications, nil
}
