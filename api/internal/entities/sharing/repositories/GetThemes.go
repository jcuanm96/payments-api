package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetThemes(ctx context.Context, req request.GetThemes) (*response.GetThemes, error) {
	query, args, squirrelErr := squirrel.Select(
		"id",
		"theme_name",
		"theme_vama_logo_color",
		"theme_icon_color",
		"theme_top_gradient_color",
		"theme_bottom_gradient_color",
		"theme_username_color",
		"theme_bio_color",
		"theme_row_color",
		"theme_row_text_color",
	).From("sharing.themes").
		Where("id > ?", req.CursorThemeID).
		OrderBy("id ASC").
		Limit(uint64(req.Limit)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		vlog.Errorf(ctx, "Error fetching themes. Err: %v", squirrelErr)
		return nil, squirrelErr
	}

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Errorf(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}
	defer rows.Close()
	themes := []response.Theme{}
	for rows.Next() {
		theme := response.Theme{}
		scanErr := rows.Scan(
			&theme.ID,
			&theme.Name,
			&theme.LogoColor,
			&theme.IconColor,
			&theme.TopGradientColor,
			&theme.BottomGradientColor,
			&theme.UsernameColor,
			&theme.BioColor,
			&theme.RowColor,
			&theme.RowTextColor,
		)

		if scanErr != nil {
			return nil, scanErr
		}

		themes = append(themes, theme)
	}
	res := response.GetThemes{
		Themes: themes,
	}
	return &res, nil
}
