package contact

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

type Repository interface {
	baserepo.BaseRepository

	CreateContact(ctx context.Context, userID, contactID int) error
	DeleteContact(ctx context.Context, userID, contactID int) error
	ValidateUserIdExists(ctx context.Context, userID int) (bool, error)

	InsertPendingContact(ctx context.Context, userID int, phone request.Phone, firstName, lastName string) error
	GetRecommendations(ctx context.Context, userID int) ([]response.User, error)
	BatchInsertContacts(ctx context.Context, currUserID int, userIDs []int, exec utils.Executable) error
	UpdatePendingContacts(ctx context.Context, phone string, userID int, exec utils.Executable) error
}
