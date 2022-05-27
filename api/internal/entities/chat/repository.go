package chat

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	MasterNode() *pgxpool.Pool
	Redis() vredis.Client

	UpdateGoatChatConversationEndTS(ctx context.Context, lastMessageCreatedAt int64, sendbirdChannelID string, runnable utils.Runnable) error
	GetPreviousConversationEndTS(ctx context.Context, channelID string) (int64, error)
	InsertGoatChat(ctx context.Context, customerUserID int, providerUserID int, req request.StartGoatChat, runnable utils.Runnable) error
	UpdateGoatChat(ctx context.Context, req request.EndGoatChat) (*UpdateGoatChatResult, error)
	GetPostConversation(ctx context.Context, postID int) (*response.Conversation, error)

	UpsertPaidGroupChat(ctx context.Context, userID int, sendbirdChannelID string, stripeProductID string, price int, currency string, link string, memberLimit int, isMemberLimitEnabled bool, metadataBytes []byte, runnable utils.Runnable) (*int, error)
	GetPaidGroup(ctx context.Context, runnable utils.Runnable, channelID string) (*response.PaidGroup, error)
	ListGoatPaidGroups(ctx context.Context, goatID int, cursorID int, limit int64) ([]response.PaidGroup, error)
	InsertBannedChatUser(ctx context.Context, bannedUserID int, userID int, channelID string, runnable utils.Runnable) error
	DeleteBannedChatUser(ctx context.Context, bannedUserID int, channelID string, runnable utils.Runnable) error
	GetBannedChatUser(ctx context.Context, runnable utils.Runnable, userID int, channelID string) (*int, error)
	ListBannedUsers(ctx context.Context, runnable utils.Runnable, req request.ListBannedUsers) (*response.ListBannedUsers, error)
	GenerateUniqueLink(ctx context.Context, groupName string) (*string, error)

	UpsertFreeGroupChat(ctx context.Context, userID int, sendbirdChannelID string, link string, memberLimit int, isMemberLimitEnabled bool, metadataBytes []byte, runnable utils.Runnable) (*int, error)
	GetFreeGroup(ctx context.Context, runnable utils.Runnable, channelID string) (*response.FreeGroup, error)
	ListUserFreeGroups(ctx context.Context, userID int, cursorID int, limit int64) ([]response.FreeGroup, error)
	DeleteFreeGroupProduct(ctx context.Context, runnable utils.Runnable, channelID string) error
	AddFreeGroupChatCoCreators(ctx context.Context, runanble utils.Runnable, req request.AddFreeGroupChatCoCreators) error
	RemoveFreeGroupChatCoCreator(ctx context.Context, runnable utils.Runnable, req request.RemoveFreeGroupChatCoCreator) error
}
