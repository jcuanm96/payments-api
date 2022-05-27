package sharing

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

type Repository interface {
	baserepo.BaseRepository
	GetBioData(ctx context.Context, runnable utils.Runnable, username string, userID int) (*response.GetBioData, error)

	GenerateMessageLinkSuffix(ctx context.Context, runnable utils.Runnable) (string, error)
	InsertMessageLink(ctx context.Context, runnable utils.Runnable, linkSuffix, messageID, channelID string) error
	GetMessageByLinkSuffix(ctx context.Context, runnable utils.Runnable, linkSuffix string) (*response.MessageInfo, error)

	GetBioLinks(ctx context.Context, userID int) (*response.BioLinks, error)
	UpsertBioLinks(ctx context.Context, userID int, textContents string, links string, themeID int) (*response.UpsertBioLinks, error)

	GetThemes(ctx context.Context, req request.GetThemes) (*response.GetThemes, error)
	UpsertTheme(ctx context.Context, req request.UpsertTheme) error
	DeleteTheme(ctx context.Context, themeID int) error
}
