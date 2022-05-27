package repositories

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetBioLinks(ctx context.Context, userID int) (*response.BioLinks, error) {
	res := response.BioLinks{}

	linksSelect := `
		SELECT
			users.first_name,
			users.last_name, 
			users.profile_avatar,
			users.username,
			bio.text_content,
			link.text_contents,
			link.links,
			link.theme_id
		FROM core.users AS users
		LEFT JOIN core.user_bio bio ON users.id = bio.user_id
		LEFT JOIN sharing.bio_links link ON users.id = link.goat_id
	  	WHERE
			users.id = $1;
	`

	bioArgs := []interface{}{userID}

	bioRow := s.MasterNode().QueryRow(ctx, linksSelect, bioArgs...)
	var themeID *int
	var textContents *string
	var links *string
	bioScanErr := bioRow.Scan(
		&res.FirstName,
		&res.LastName,
		&res.ProfileAvatar,
		&res.Username,
		&res.BioText,
		&textContents,
		&links,
		&themeID,
	)
	if bioScanErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(bioScanErr, linksSelect, bioArgs))
		return nil, bioScanErr
	}
	if textContents != nil {
		res.TextContents = strings.Split(*textContents, "|")
	}
	if links != nil {
		res.Links = strings.Split(*links, "|")
	}

	queryBuilder := squirrel.Select(
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
	).
		From("sharing.themes")

	if themeID != nil {
		queryBuilder = queryBuilder.Where("id = ?", *themeID)
	}

	themeQuery, themeArgs, squirrelErr := queryBuilder.
		OrderBy("id ASC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	themeRow := s.MasterNode().QueryRow(ctx, themeQuery, themeArgs...)
	res.Theme = response.Theme{}
	themeScanErr := themeRow.Scan(
		&res.Theme.ID,
		&res.Theme.Name,
		&res.Theme.LogoColor,
		&res.Theme.IconColor,
		&res.Theme.TopGradientColor,
		&res.Theme.BottomGradientColor,
		&res.Theme.UsernameColor,
		&res.Theme.BioColor,
		&res.Theme.RowColor,
		&res.Theme.RowTextColor,
	)
	if themeScanErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(themeScanErr, themeQuery, themeArgs))
		return nil, themeScanErr
	}

	return &res, nil
}
