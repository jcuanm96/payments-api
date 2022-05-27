package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jmoiron/sqlx/types"
)

func UpsertPaidGroupChat(
	ctx context.Context,
	userID int,
	sendbirdChannelID string,
	stripeProductID string,
	priceInSmallestDenom int,
	currency string,
	link string,
	memberLimit int,
	isMemberLimitEnabled bool,
	metadataBytes []byte,
	runnable utils.Runnable,
) (*int, error) {
	query, args, squirrelErr := squirrel.Insert("product.paid_group_chats").
		Columns(
			"price_in_smallest_denom",
			"currency",
			"goat_user_id",
			"sendbird_channel_id",
			"stripe_product_id",
			"link_suffix",
			"metadata",
			"member_limit",
			"is_member_limit_enabled",
		).
		Values(
			priceInSmallestDenom,
			currency,
			userID,
			sendbirdChannelID,
			stripeProductID,
			link,
			types.JSONText(metadataBytes),
			memberLimit,
			isMemberLimitEnabled,
		).
		Suffix(`
		ON CONFLICT (sendbird_channel_id)
		DO UPDATE SET 
		price_in_smallest_denom = ?,
		currency = ?,
		link_suffix = ?,
		metadata = ?,
		member_limit = ?,
		is_member_limit_enabled = ?`,
			priceInSmallestDenom,
			currency,
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
		vlog.Errorf(ctx, "Error constructing paid group chat insert query: %s", squirrelErr.Error())
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

func (s *repository) UpsertPaidGroupChat(
	ctx context.Context,
	userID int,
	sendbirdChannelID string,
	stripeProductID string,
	priceInSmallestDenom int,
	currency string,
	link string,
	memberLimit int,
	isMemberLimitEnabled bool,
	metadataBytes []byte,
	runnable utils.Runnable,
) (*int, error) {
	return UpsertPaidGroupChat(
		ctx,
		userID,
		sendbirdChannelID,
		stripeProductID,
		priceInSmallestDenom,
		currency,
		link,
		memberLimit,
		isMemberLimitEnabled,
		metadataBytes,
		runnable,
	)
}
