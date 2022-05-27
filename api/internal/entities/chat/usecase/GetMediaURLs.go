package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/google/uuid"
)

func (svc *usecase) GetMediaURLs(ctx context.Context) (*response.MediaURLs, error) {
	object := uuid.NewString()
	expires := time.Now().Add(10 * time.Minute)

	var wg sync.WaitGroup

	var url string
	var signedURLErr error
	var downloadURL string
	wg.Add(1)
	go func() {
		defer wg.Done()
		url, signedURLErr = svc.gcsClient.ChatMediaBucket.SignedURL(object, &storage.SignedURLOptions{
			Method:  http.MethodPut,
			Expires: expires,
		})
		downloadURL = fmt.Sprintf(constants.GCS_URL_F, svc.gcsClient.ChatMediaBucketName, object)
	}()

	var thumbnailURL string
	var thumbnailSignedURLErr error
	var thumbnailDownloadURL string
	wg.Add(1)
	go func() {
		defer wg.Done()
		thumbnailObject := fmt.Sprintf("%s-thumbnail", object)
		thumbnailURL, thumbnailSignedURLErr = svc.gcsClient.ChatMediaBucket.SignedURL(thumbnailObject, &storage.SignedURLOptions{
			Method:  http.MethodPut,
			Expires: expires,
		})
		thumbnailDownloadURL = fmt.Sprintf(constants.GCS_URL_F, svc.gcsClient.ChatMediaBucketName, thumbnailObject)
	}()

	wg.Wait()

	if signedURLErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error generating signedURL for media upload: %v", signedURLErr),
		)
	}

	if thumbnailSignedURLErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error generating thumbnailSignedURL for media upload: %v", thumbnailSignedURLErr),
		)
	}

	res := &response.MediaURLs{
		UploadURL:    url,
		ThumbnailURL: thumbnailURL,
		ObjectID:     object,

		DownloadURL:          downloadURL,
		ThumbnailDownloadURL: thumbnailDownloadURL,
	}
	return res, nil
}
