package upload

import (
	"context"
	"fmt"
	"io"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/unique"
)

type Mock struct {
}

func (m Mock) ProfileAvatarFileName(user *response.User, ext string) string {
	return fmt.Sprintf("profile-avatar-%s-%s-%s%s", user.FirstName, user.LastName, unique.New(12), ext)
}

func (m Mock) Upload(ctx context.Context, filename string, r io.Reader) error {
	return nil
}
