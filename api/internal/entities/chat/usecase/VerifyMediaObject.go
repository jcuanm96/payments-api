package service

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) VerifyMediaObject(ctx context.Context, req request.VerifyMediaObject) (*response.VerifyMediaObject, error) {
	mediaObjectHandle := svc.gcsClient.ChatMediaBucket.Object(req.ObjectID)
	_, attrErr := mediaObjectHandle.Attrs(ctx)
	if attrErr == storage.ErrObjectNotExist {
		return &response.VerifyMediaObject{Exists: false}, nil
	} else if attrErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error something went wrong when verifying media object: %v", attrErr),
		)
	}

	return &response.VerifyMediaObject{Exists: true}, nil
}
