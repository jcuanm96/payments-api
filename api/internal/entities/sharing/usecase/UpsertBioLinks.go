package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const upsertBioLinksErr = "Something went wrong when trying to edit bio links."

func (svc *usecase) UpsertBioLinks(ctx context.Context, req request.UpsertBioLinks) (*response.UpsertBioLinks, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			upsertBioLinksErr,
			fmt.Sprintf("Could not find user in the current context: %v.", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			upsertBioLinksErr,
			"user returned nil for UpsertBioLinks",
		)
	}
	textContents := strings.Join(req.TextContents, "|")
	links := strings.Join(req.Links, "|")

	res, getBioLinksErr := svc.repo.UpsertBioLinks(ctx, user.ID, textContents, links, req.ThemeID)
	if getBioLinksErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			upsertBioLinksErr,
			fmt.Sprintf("Could not get bio data for user %d. Err: %v.", user.ID, getBioLinksErr),
		)
	}

	return res, nil
}
