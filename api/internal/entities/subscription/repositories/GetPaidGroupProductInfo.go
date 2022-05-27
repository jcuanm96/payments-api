package repositories

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/subscription"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (s *repository) GetPaidGroupProductInfo(ctx context.Context, runnable utils.Runnable, channelID string) (*subscription.PaidGroupChatInfo, error) {
	return GetPaidGroupProductInfo(ctx, runnable, channelID)
}

func GetPaidGroupProductInfo(ctx context.Context, runnable utils.Runnable, channelID string) (*subscription.PaidGroupChatInfo, error) {
	query, args, squirrelErr := squirrel.Select(
		"price_in_smallest_denom",
		"currency",
		"goat_user_id",
		"sendbird_channel_id",
		"stripe_product_id",
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

	res := &subscription.PaidGroupChatInfo{}
	var metadataBytes []byte
	scanErr := row.Scan(
		&res.PriceInSmallestDenom,
		&res.Currency,
		&res.GoatID,
		&res.ChannelID,
		&res.StripeProductID,
		&res.LinkSuffix,
		&res.MemberLimit,
		&res.IsMemberLimitEnabled,
		&metadataBytes,
	)

	if scanErr != nil {
		return nil, scanErr
	}

	metadata := response.GroupMetadata{}
	unmarshalErr := json.Unmarshal(metadataBytes, &metadata)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	res.Metadata = metadata

	return res, nil
}
