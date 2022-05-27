package auth

import (
	"context"

	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	MasterNode() *pgxpool.Pool
	baserepo.BaseCacheRepository
	baserepo.BaseRepository
	InsertUserInviteCodes(ctx context.Context, exec utils.Executable, userID int, codes []string) error
	GetPendingContactsByPhone(ctx context.Context, phone string) ([]int, error)
}
