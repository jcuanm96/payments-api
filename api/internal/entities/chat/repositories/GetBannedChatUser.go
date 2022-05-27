package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetBannedChatUser(ctx context.Context, runnable utils.Runnable, userID int, channelID string) (*int, error) {
	query, args, squirrelErr := squirrel.Select(
		"banned_user_id",
	).
		From("core.banned_chat_users").
		Where("banned_user_id = ?", userID).
		Where("sendbird_channel_id = ?", channelID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	row := runnable.QueryRow(ctx, query, args...)

	var bannedUserID int
	scanErr := row.Scan(
		&bannedUserID,
	)
	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, scanErr
	}

	return &bannedUserID, nil
}
