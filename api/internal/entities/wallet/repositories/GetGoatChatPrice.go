package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetGoatChatPrice(ctx context.Context, goatUserID int) (*response.GetGoatChatPrice, error) {
	var res response.GetGoatChatPrice
	query, args, squirrelErr := squirrel.Select(
		"price_in_smallest_denom",
		"currency",
	).
		From("product.goat_chats").
		Where("goat_user_id = ?", goatUserID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if squirrelErr != nil {
		return nil, squirrelErr
	}

	row := s.MasterNode().QueryRow(ctx, query, args...)
	scanErr := row.Scan(
		&res.PriceInSmallestDenom,
		&res.Currency,
	)
	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return nil, nil
		}
		vlog.Error(ctx, utils.SqlErrLogMsg(scanErr, query, args))
		return nil, scanErr
	}
	return &res, nil
}
