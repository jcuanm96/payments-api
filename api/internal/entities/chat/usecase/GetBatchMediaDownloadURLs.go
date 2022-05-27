package service

import (
	"context"
	"fmt"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func (svc *usecase) GetBatchMediaDownloadURLs(ctx context.Context, req request.GetBatchMediaDownloadURLs) (*response.BatchMediaDownloadURLs, error) {
	urls := []response.MediaDownloadURL{}

	// expire in 200 years (never expire)
	expiresAt := time.Now().AddDate(200, 0, 0)

	for _, objectID := range req.ObjectIDs {
		downloadURL := fmt.Sprintf(constants.GCS_URL_F, svc.gcsClient.ChatMediaBucketName, objectID)

		thumbnailObjectID := fmt.Sprintf("%s-thumbnail", objectID)
		thumbnailDownloadURL := fmt.Sprintf(constants.GCS_URL_F, svc.gcsClient.ChatMediaBucketName, thumbnailObjectID)

		url := response.MediaDownloadURL{
			ObjectID:     objectID,
			DownloadURL:  downloadURL,
			ThumbnailURL: thumbnailDownloadURL,
			ExpiresAt:    expiresAt,
		}

		urls = append(urls, url)
	}

	res := &response.BatchMediaDownloadURLs{
		URLs: urls,
	}
	return res, nil
}
