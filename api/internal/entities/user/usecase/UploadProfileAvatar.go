package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errUploadingProfilePicture = "Something went wrong uploading your profile picture. Please try again."

func (svc *usecase) UploadProfileAvatar(ctx context.Context, req request.UploadProfileAvatar) (*response.UploadProfileAvatar, error) {
	user, getCurrentUserErr := svc.GetCurrentUser(ctx)
	if getCurrentUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUploadingProfilePicture,
			fmt.Sprintf("Error occurred when retrieving user from db: %v", getCurrentUserErr),
		)
	} else if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errUploadingProfilePicture,
			"user was nil in UploadProfileAvatar",
		)
	}

	file, openFileErr := req.ProfileAvatar.Open()
	if openFileErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUploadingProfilePicture,
			fmt.Sprintf("Error occurred when opening the profile avatar file: %v", openFileErr),
		)
	}

	defer file.Close()

	filename := fmt.Sprintf("profile-avatar-%d%s", user.ID, filepath.Ext(req.ProfileAvatar.Filename))

	writerCtx := svc.gcsClient.ProfileAvatarBucket.Object(filename).NewWriter(ctx)
	if _, copyFileErr := io.Copy(writerCtx, file); copyFileErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUploadingProfilePicture,
			fmt.Sprintf("Error occurred when uploading file to GCS for user %d: %v", user.ID, copyFileErr),
		)
	}
	if closeFileErr := writerCtx.Close(); closeFileErr != nil {
		vlog.Errorf(ctx, "Error occurred when closing file to GCS for user %d: %v", user.ID, closeFileErr)
	}
	fileUrl := fmt.Sprintf(constants.GCS_URL_F, svc.gcsClient.ProfileAvatarBucketName, filename)

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUploadingProfilePicture,
			fmt.Sprintf("Could not begin transaction: %v", txErr),
		)
	}
	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	updateAvatarErr := svc.repo.UpdateUserProfileAvatar(ctx, tx, user.ID, fileUrl)
	if updateAvatarErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUploadingProfilePicture,
			fmt.Sprintf("Error occurred when updating user profile avatar to db for user %d and file url %s: %v", user.ID, fileUrl, updateAvatarErr),
		)
	}

	updateSendbirdUserParams := &sendbird.UpdateUserParams{}
	updateSendbirdUserParams.ProfileURL = fileUrl
	_, updateUserErr := svc.sendBirdClient.UpdateUser(user.ID, updateSendbirdUserParams)

	if updateUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUploadingProfilePicture,
			fmt.Sprintf("Something went wrong when updating SendBird user for uuid %s: %v", user.UUID, updateUserErr),
		)
	}

	res := response.UploadProfileAvatar{
		Filename: fileUrl,
	}

	commit = true
	return &res, nil
}
