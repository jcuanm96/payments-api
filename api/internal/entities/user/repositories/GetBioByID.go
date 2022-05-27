package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetBioByID(ctx context.Context, userID int) (string, error) {
	query, args, squirrelErr := squirrel.Select(
		"text_content",
	).
		From("core.user_bio").
		Where("user_id = ?", userID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return "", squirrelErr
	}

	row, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return "", queryErr
	}

	defer row.Close()
	// It's valid for a user to not have a bio
	if !row.Next() {
		return "", nil
	}

	var bio string
	scanErr := row.Scan(&bio)
	if scanErr != nil {
		return "", nil
	}

	return bio, nil
}
