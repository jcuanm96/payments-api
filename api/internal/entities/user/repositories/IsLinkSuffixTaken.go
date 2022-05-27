package repositories

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4"
)

func (s *repository) IsLinkSuffixTaken(ctx context.Context, runnable utils.Runnable, suffix string) (*user.LinkSuffixTaken, error) {
	return IsLinkSuffixTaken(ctx, runnable, suffix)
}

func IsLinkSuffixTaken(ctx context.Context, runnable utils.Runnable, suffix string) (*user.LinkSuffixTaken, error) {
	// Username Check
	usernameQuery, usernameArgs, usernameSquirrelErr := squirrel.Select(
		"id",
	).
		From("core.users").
		Where("LOWER(username) = ?", strings.ToLower(suffix)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if usernameSquirrelErr != nil {
		return nil, usernameSquirrelErr
	}

	usernameRow := runnable.QueryRow(ctx, usernameQuery, usernameArgs...)

	var userID int
	usernameScanErr := usernameRow.Scan(&userID)

	if usernameScanErr == nil {
		res := &user.LinkSuffixTaken{
			IsTaken: true,
			UserID:  &userID,
		}
		return res, nil
	} else if usernameScanErr != pgx.ErrNoRows {
		return nil, usernameScanErr
	}

	// Paid Group Check
	linkQuery, linkArgs, linkSquirrelErr := squirrel.Select(
		"sendbird_channel_id",
	).
		From("product.paid_group_chats").
		Where("LOWER(link_suffix) = ?", strings.ToLower(suffix)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if linkSquirrelErr != nil {
		return nil, linkSquirrelErr
	}

	linkRow := runnable.QueryRow(ctx, linkQuery, linkArgs...)

	var paidChannelID string
	linkScanErr := linkRow.Scan(&paidChannelID)
	if linkScanErr == nil {
		res := &user.LinkSuffixTaken{
			IsTaken:            true,
			PaidGroupChannelID: &paidChannelID,
		}
		return res, nil
	} else if linkScanErr != nil && linkScanErr != pgx.ErrNoRows {
		return nil, linkScanErr
	}

	// Free Group
	freeGroupLinkQuery, freeGroupLinkArgs, freeGroupLinkSquirrelErr := squirrel.Select(
		"sendbird_channel_id",
	).
		From("product.free_group_chats").
		Where("LOWER(link_suffix) = ?", strings.ToLower(suffix)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if freeGroupLinkSquirrelErr != nil {
		return nil, freeGroupLinkSquirrelErr
	}

	freeGroupLinkRow := runnable.QueryRow(ctx, freeGroupLinkQuery, freeGroupLinkArgs...)

	var freeChannelID string
	freeGroupLinkScanErr := freeGroupLinkRow.Scan(&freeChannelID)
	if freeGroupLinkScanErr == nil {
		res := &user.LinkSuffixTaken{
			IsTaken:            true,
			FreeGroupChannelID: &freeChannelID,
		}
		return res, nil
	} else if freeGroupLinkScanErr != nil && freeGroupLinkScanErr != pgx.ErrNoRows {
		return nil, freeGroupLinkScanErr
	}

	res := &user.LinkSuffixTaken{
		IsTaken: false,
	}

	return res, nil
}
