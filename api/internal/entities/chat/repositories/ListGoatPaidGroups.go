package repositories

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) ListGoatPaidGroups(ctx context.Context, goatID int, cursorID int, limit int64) ([]response.PaidGroup, error) {
	queryBuilder := squirrel.Select(
		"id",
		"price_in_smallest_denom",
		"currency",
		"sendbird_channel_id",
		"link_suffix",
		"metadata",
	).
		From("product.paid_group_chats").
		Where("goat_user_id = ?", goatID)

	if cursorID > 0 {
		queryBuilder = queryBuilder.Where("id < ?", cursorID)
	}

	query, args, squirrelErr := queryBuilder.
		OrderBy("id DESC").
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}
	defer rows.Close()

	groups := []response.PaidGroup{}
	for rows.Next() {
		group := response.PaidGroup{}
		var metadataBytes []byte
		scanErr := rows.Scan(
			&group.ID,
			&group.PriceInSmallestDenom,
			&group.Currency,
			&group.ChannelID,
			&group.LinkSuffix,
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

		group.Metadata = metadata

		groups = append(groups, group)
	}

	return groups, nil
}
