package contact

import (
	"context"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

type Usecase interface {
	GetContacts(ctx context.Context, req request.GetContacts) (*response.GetContacts, error)
	IsContact(ctx context.Context, contactID int) (*response.IsContact, error)
	CreateContact(ctx context.Context, contactID int) (*response.User, error)
	DeleteContact(ctx context.Context, contactID int) (*response.DeleteContact, error)

	UploadContacts(ctx context.Context, req request.UploadContacts) error
	GetRecommendations(ctx context.Context, limit int) (*response.ContactRecommendations, error)
	BatchAddUsersToContacts(ctx context.Context, req cloudtasks.AddUserContactsTask) error
}
