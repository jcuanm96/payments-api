package repositories

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) ListUserFreeGroups(ctx context.Context, userID int, cursorID int, limit int64) ([]response.FreeGroup, error) {
	queryBuilder := squirrel.Select(
		"id",
		"sendbird_channel_id",
		"link_suffix",
		"metadata",
	).
		From("product.free_group_chats").
		Where("creator_user_id = ?", userID)

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

	groups := []response.FreeGroup{}
	for rows.Next() {
		group := response.FreeGroup{}
		var metadataBytes []byte
		scanErr := rows.Scan(
			&group.ID,
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
