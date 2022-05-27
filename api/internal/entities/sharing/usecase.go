package sharing

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

type Usecase interface {
	GetBioData(ctx context.Context, req request.GetBioData) (*response.GetBioData, error)
	UpsertBioLinks(ctx context.Context, req request.UpsertBioLinks) (*response.UpsertBioLinks, error)
	GetLink(ctx context.Context, req request.GetLink) (*response.GetVamaMeLink, error)

	NewMessageLink(ctx context.Context, req request.NewMessageLink) (*response.MessageLink, error)
	PublicGetMessageByLink(ctx context.Context, req request.GetMessageByLink) (*response.PublicMessage, error)
	GetMessageByLinkSuffix(ctx context.Context, req request.GetMessageByLink) (*response.MessageInfo, error)

	GetThemes(ctx context.Context, req request.GetThemes) (*response.GetThemes, error)
	UpsertTheme(ctx context.Context, req request.UpsertTheme) (*response.UpsertTheme, error)
	DeleteTheme(ctx context.Context, req request.DeleteTheme) error
}
