package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) UpsertTheme(ctx context.Context, req request.UpsertTheme) error {
	query, args, squirrelErr := squirrel.Insert("sharing.themes").
		Columns(
			"theme_name",
			"theme_vama_logo_color",
			"theme_icon_color",
			"theme_top_gradient_color",
			"theme_bottom_gradient_color",
			"theme_username_color",
			"theme_bio_color",
			"theme_row_color",
			"theme_row_text_color",
		).
		Values(
			req.Name,
			req.LogoColor,
			req.IconColor,
			req.TopGradientColor,
			req.BottomGradientColor,
			req.UsernameColor,
			req.BioColor,
			req.RowColor,
			req.RowTextColor,
		).
		Suffix(`
		ON CONFLICT (theme_name)
		DO UPDATE SET
			theme_vama_logo_color = ?,
			theme_icon_color = ?,
			theme_top_gradient_color = ?,
			theme_bottom_gradient_color = ?,
			theme_username_color = ?,
			theme_bio_color = ?,
			theme_row_color = ?,
			theme_row_text_color = ?`,
			req.LogoColor,
			req.IconColor,
			req.TopGradientColor,
			req.BottomGradientColor,
			req.UsernameColor,
			req.BioColor,
			req.RowColor,
			req.RowTextColor,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		vlog.Errorf(ctx, "Error generating SQL query for Upsert Bio Links. Err: %s", squirrelErr.Error())
		return squirrelErr
	}

	_, queryErr := s.MasterNode().Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
