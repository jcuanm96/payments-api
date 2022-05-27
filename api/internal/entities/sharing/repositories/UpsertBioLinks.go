package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) UpsertBioLinks(ctx context.Context, userID int, textContents string, links string, themeID int) (*response.UpsertBioLinks, error) {
	query, args, squirrelErr := squirrel.Insert("sharing.bio_links").
		Columns(
			"goat_id",
			"text_contents",
			"links",
			"theme_id",
		).
		Values(
			userID,
			textContents,
			links,
			themeID,
		).
		Suffix(`
		ON CONFLICT (goat_id)
		DO UPDATE SET
			text_contents = ?,
			links = ?,
			theme_id = ?`,
			textContents,
			links,
			themeID,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		vlog.Errorf(ctx, "Error generating SQL query for Upsert Bio Links. Err: %s", squirrelErr.Error())
		return nil, squirrelErr
	}

	_, queryErr := s.MasterNode().Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}

	return &response.UpsertBioLinks{}, nil
}
