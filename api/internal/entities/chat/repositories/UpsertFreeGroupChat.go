package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jmoiron/sqlx/types"
)

func UpsertFreeGroupChat(
	ctx context.Context,
	userID int,
	sendbirdChannelID string,
	link string,
	memberLimit int,
	isMemberLimitEnabled bool,
	metadataBytes []byte,
	runnable utils.Runnable,
) (*int, error) {
	query, args, squirrelErr := squirrel.Insert("product.free_group_chats").
		Columns(
			"creator_user_id",
			"sendbird_channel_id",
			"link_suffix",
			"metadata",
			"member_limit",
			"is_member_limit_enabled",
		).
		Values(
			userID,
			sendbirdChannelID,
			link,
			types.JSONText(metadataBytes),
			memberLimit,
			isMemberLimitEnabled,
		).
		Suffix(`
		ON CONFLICT (sendbird_channel_id)
		DO UPDATE SET 
		link_suffix = ?,
		metadata = ?,
		member_limit = ?,
		is_member_limit_enabled = ?`,
			link,
			types.JSONText(metadataBytes),
			memberLimit,
			isMemberLimitEnabled,
		).
		Suffix(
			`RETURNING id`,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}
	rows := runnable.QueryRow(ctx, query, args...)
	var id int
	scanErr := rows.Scan(
		&id,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	return &id, nil
}

func (s *repository) UpsertFreeGroupChat(
	ctx context.Context,
	userID int,
	sendbirdChannelID string,
	link string,
	memberLimit int,
	isMemberLimitEnabled bool,
	metadataBytes []byte,
	runnable utils.Runnable,
) (*int, error) {
	return UpsertFreeGroupChat(
		ctx,
		userID,
		sendbirdChannelID,
		link,
		memberLimit,
		isMemberLimitEnabled,
		metadataBytes,
		runnable,
	)
}
