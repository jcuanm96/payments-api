package repositories

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetPaidGroup(ctx context.Context, runnable utils.Runnable, channelID string) (*response.PaidGroup, error) {
	return GetPaidGroup(ctx, runnable, channelID)
}

func GetPaidGroup(ctx context.Context, runnable utils.Runnable, channelID string) (*response.PaidGroup, error) {
	query, args, squirrelErr := squirrel.Select(
		"id",
		"goat_user_id",
		"price_in_smallest_denom",
		"currency",
		"sendbird_channel_id",
		"link_suffix",
		"member_limit",
		"is_member_limit_enabled",
		"metadata",
	).
		From("product.paid_group_chats").
		Where("sendbird_channel_id = ?", channelID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	row := runnable.QueryRow(ctx, query, args...)

	group := &response.PaidGroup{}
	var metadataBytes []byte
	scanErr := row.Scan(
		&group.ID,
		&group.GoatID,
		&group.PriceInSmallestDenom,
		&group.Currency,
		&group.ChannelID,
		&group.LinkSuffix,
		&group.MemberLimit,
		&group.IsMemberLimitEnabled,
		&metadataBytes,
	)

	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, scanErr
	}

	metadata := response.GroupMetadata{}
	unmarshalErr := json.Unmarshal(metadataBytes, &metadata)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	group.Metadata = metadata
	return group, nil
}
