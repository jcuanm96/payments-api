package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GenerateMessageLinkSuffix(ctx context.Context, runnable utils.Runnable) (string, error) {
	const attempts = 6
	const linkSuffixLength = 10

	for i := 0; i < attempts; i++ {
		linkSuffix := utils.RandAlphaNumeric(linkSuffixLength)
		taken, checkTakenErr := s.checkMessageLinkSuffixTaken(ctx, runnable, linkSuffix)
		if checkTakenErr != nil {
			return "", checkTakenErr
		}
		if !taken {
			return linkSuffix, nil
		}

	}
	errMsg := fmt.Sprintf("unable to generate a unique message link suffix after %d attempts", attempts)
	telegram.TelegramClient.SendMessage(errMsg)
	return "", errors.New(errMsg)
}

func (s *repository) checkMessageLinkSuffixTaken(ctx context.Context, runnable utils.Runnable, linkSuffix string) (bool, error) {
	query, args, squirrelErr := squirrel.Select("id").
		From("sharing.message_links").
		Where("link_suffix = ?", linkSuffix).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if squirrelErr != nil {
		return true, squirrelErr
	}

	row := runnable.QueryRow(ctx, query, args...)
	var id int
	scanErr := row.Scan(&id)
	if scanErr == pgx.ErrNoRows {
		return false, nil
	}

	return true, scanErr
}
