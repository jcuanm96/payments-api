package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (s *repository) GetMessageByLinkSuffix(ctx context.Context, runnable utils.Runnable, linkSuffix string) (*response.MessageInfo, error) {
	query, args, squirrelErr := squirrel.Select(
		"sendbird_channel_id",
		"message_id",
	).
		From("sharing.message_links").
		Where("link_suffix = ?", linkSuffix).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	row := runnable.QueryRow(ctx, query, args...)

	var info response.MessageInfo
	scanErr := row.Scan(
		&info.ChannelID,
		&info.MessageID,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	return &info, nil
}
