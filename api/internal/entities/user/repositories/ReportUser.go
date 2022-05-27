package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) ReportUser(ctx context.Context, reporterUserID int, reportedUserID int, description string) error {
	query, args, squirrelErr := squirrel.Insert("core.user_reports").
		Columns(
			"reporter_user_id",
			"reported_user_id",
			"description",
		).
		Values(
			reporterUserID,
			reportedUserID,
			description,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if squirrelErr != nil {
		return squirrelErr
	}

	_, queryErr := s.MasterNode().Exec(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return queryErr
	}

	return nil
}
