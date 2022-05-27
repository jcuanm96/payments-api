package search

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

type Usecase interface {
	SearchGlobal(ctx context.Context, req request.SearchGlobal) (*response.SearchGlobal, error)
	SearchMention(ctx context.Context, req request.SearchMention) (*response.SearchMention, error)
	SearchUsers(ctx context.Context, userID int, req request.NewGridList) (*response.Paging, []response.User, error)
	GetSendbirdChannelsOwnedBy(ctx context.Context, userID int, req request.SearchGlobal) ([]sendbird.GroupChannel, error)
	GetMyGroupSendbirdChannels(userID int, req request.SearchGlobal) ([]sendbird.GroupChannel, error)
	GetMyDistinctSendbirdChannels(userID int, req request.SearchGlobal) ([]sendbird.GroupChannel, error)
	GetGlobalSendbirdChannels(userID int, req request.SearchGlobal) ([]sendbird.GroupChannel, error)
	GetChannelsByLink(ctx context.Context, userID int, req request.SearchGlobal) (*Channels, error)
	GetAllSendbirdChannels(ctx context.Context, userID int, req request.SearchGlobal) (*Channels, error)
	GetChannels(ctx context.Context, userID int, req request.SearchGlobal) (*Channels, error)
	GetUsers(ctx context.Context, query string, userID int, limit int) ([]response.User, error)
	GetMessages(ctx context.Context, query string, userID int, limit int) ([]response.SearchMessage, error)
}
