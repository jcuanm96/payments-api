package service

import (
	"context"
	"fmt"
	imagelib "image"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/upload"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/google/uuid"
)

const errMakingFeedPost = "Something went wrong when making feed post."

func (svc *usecase) MakeFeedPost(ctx context.Context, req request.MakeFeedPost) (*response.FeedPost, error) {
	var wg sync.WaitGroup
	var linkSuffix string
	var generateLinkSuffixErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		linkSuffix, generateLinkSuffixErr = svc.repo.GenerateFeedPostLinkSuffix(ctx)
	}()

	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errMakingFeedPost,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errMakingFeedPost,
			"user was nil in MakeFeedPost",
		)
	} else if user.Type != "GOAT" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Only creators can make feed posts. Interested?  Sign up to be a creator now!",
			fmt.Sprintf("user was of type %s", user.Type),
		)
	}

	var image response.PostImage
	if req.Image != nil {
		newImage, imageErr := uploadFeedPostImage(ctx, svc.gcsClient, req.Image, user.ID)
		if imageErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				errMakingFeedPost,
				fmt.Sprintf("Error uploading feed post image: %v", imageErr),
			)
		}
		image = newImage
	}

	wg.Wait()
	if generateLinkSuffixErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errMakingFeedPost,
			fmt.Sprintf("Error generating link suffix: %v", generateLinkSuffixErr),
		)
	}

	newPostID, sqlErr := svc.repo.UpsertFeedPost(ctx, req, user.ID, image, linkSuffix)
	if sqlErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errMakingFeedPost,
			fmt.Sprintf("Could not make post for user %d. Err: %v", user.ID, sqlErr),
		)
	}

	post, getPostErr := svc.GetFeedPostByID(ctx, request.GetFeedPostByID{PostID: newPostID})
	if getPostErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errMakingFeedPost,
			fmt.Sprintf("Could not get new post with id %d for user %d. Err: %v", newPostID, user.ID, getPostErr),
		)
	}

	if post == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errMakingFeedPost,
			fmt.Sprintf("Post with id %d was nil for user %d", newPostID, user.ID),
		)
	}

	notificationErr := (*svc.push).SendGoatPostNotification(ctx, user.ID, post.ID)
	if notificationErr != nil {
		vlog.Errorf(ctx, "Error sending creator %d post %d notification: %v", user.ID, post.ID, notificationErr)
	}

	return post, nil
}

func uploadFeedPostImage(ctx context.Context, gcsClient *upload.Client, reqImage *multipart.FileHeader, userID int) (response.PostImage, error) {
	image := response.PostImage{}
	file, imageErr := reqImage.Open()
	if imageErr != nil {
		vlog.Errorf(ctx, "Error occurred when opening the post feed file for user %d: %v", userID, imageErr)
		return image, imageErr
	}

	filename := fmt.Sprintf("feed-post-user-%d-%s%s", userID, uuid.New(), filepath.Ext(reqImage.Filename))

	writerCtx := gcsClient.FeedPostBucket.Object(filename).NewWriter(ctx)
	if _, copyErr := io.Copy(writerCtx, file); copyErr != nil {
		vlog.Errorf(ctx, "Error occurred when uploading feed post image to GCS for user %d. Err: %v", userID, copyErr)
		return image, copyErr
	}
	if writerErr := writerCtx.Close(); writerErr != nil {
		vlog.Errorf(ctx, "Error occurred when closing file to GCS for user %d. Err: %v", userID, writerErr)
	}

	// Get image dimensions
	// Reset offset in image file to 0 after being read above
	_, seekErr := file.Seek(0, 0)
	if seekErr != nil {
		vlog.Errorf(ctx, "Error seeking back to offset 0 for image file. Err: %v", seekErr)
		return image, seekErr
	}
	imageConfig, _, imageDecodeErr := imagelib.DecodeConfig(file) // Image Config Struct
	if imageDecodeErr != nil {
		vlog.Errorf(ctx, "Error occurred when decoding image file into image config. Err: %v", imageDecodeErr)
		return image, imageDecodeErr
	}

	image.URL = fmt.Sprintf(constants.GCS_URL_F, gcsClient.FeedPostBucketName, filename)
	image.Width = imageConfig.Width
	image.Height = imageConfig.Height
	return image, nil
}
