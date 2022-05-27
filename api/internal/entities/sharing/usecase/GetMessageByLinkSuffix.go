package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/jackc/pgx/v4"
)

func (svc *usecase) GetMessageByLinkSuffix(ctx context.Context, req request.GetMessageByLink) (*response.MessageInfo, error) {
	messageInfo, getMessageErr := svc.repo.GetMessageByLinkSuffix(ctx, svc.repo.MasterNode(), req.LinkSuffix)
	if getMessageErr == pgx.ErrNoRows {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"We couldn't find anything for this link.",
			"We couldn't find anything for this link.",
		)
	} else if getMessageErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprint(ctx, "Error getting message for suffix %s: %v", req.LinkSuffix, getMessageErr),
		)
	}

	return messageInfo, nil
}
