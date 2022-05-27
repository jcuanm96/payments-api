package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	cloudTasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errDeletingPaidGroup = "Something went wrong trying to delete paid group."

func (svc *usecase) DeletePaidGroup(ctx context.Context, req request.DeletePaidGroup) error {
	currUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingPaidGroup,
			fmt.Sprintf("Error getting user from current context: %v", getCurrUserErr),
		)
	} else if currUser == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errDeletingPaidGroup,
			"Could not find user in the current context",
		)
	}

	getGroupChannelParams := sendbird.GetGroupChannelParams{
		ShowMember: true,
	}
	channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(req.ChannelID, getGroupChannelParams)
	if getChannelErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingPaidGroup,
			fmt.Sprintf("Error getting Sendbird channel for paid group chat %s and user %d: %v", req.ChannelID, currUser.ID, getChannelErr),
		)
	}

	currUserIDStr := fmt.Sprint(currUser.ID)
	if !channel.HasOperator(currUserIDStr) {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You cannot delete this paid group chat.",
			"You cannot delete this paid group chat.",
		)
	}

	batchSize := 10
	userIDsBatch := []int{}
	for _, member := range channel.Members {
		// The creator shouldn't be subscribed to themselves
		if member.UserID == currUserIDStr {
			continue
		}

		userID, strconvErr := strconv.Atoi(member.UserID)
		if strconvErr != nil {
			vlog.Errorf(ctx, "Failed to convert %s to int DeletePaidGroup. Err: %v. Batch: %v", member.UserID, strconvErr, userIDsBatch)
			continue
		}

		userIDsBatch = append(userIDsBatch, userID)
		currNumIDs := len(userIDsBatch)

		if currNumIDs >= batchSize {
			sendTaskBatchErr := svc.sendTaskBatch(ctx, userIDsBatch, req.ChannelID, currUser.ID)
			if sendTaskBatchErr != nil {
				vlog.Errorf(ctx, "Error creating task when trying to delete paid group chat %s by user %d. Err: %v. Batch: %v", req.ChannelID, currUser.ID, sendTaskBatchErr, userIDsBatch)
				continue
			}

			userIDsBatch = []int{}
		}
	}

	if len(userIDsBatch) > 0 {
		sendTaskBatchErr := svc.sendTaskBatch(ctx, userIDsBatch, req.ChannelID, currUser.ID)
		if sendTaskBatchErr != nil {
			vlog.Errorf(ctx, "Error creating task as final element when trying to delete paid group chat %s by user %d: %v. Batch: %v", req.ChannelID, currUser.ID, sendTaskBatchErr, userIDsBatch)
		}
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingPaidGroup,
			fmt.Sprintf("Error begining transaction in DeletePaidGroup: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	deleteProductErr := svc.subscription.DeletePaidGroupProduct(ctx, tx, req.ChannelID)
	if deleteProductErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingPaidGroup,
			fmt.Sprintf("Error deleting paid group chat product for channel %s and user %d: %v", req.ChannelID, currUser.ID, deleteProductErr),
		)
	}

	// Deleting Channel after Cloud Task batching succeeds, because worst case scenario, the creator
	// can retry to delete the channel and we'll just get duplicate events going into target endpoint
	deleteChannelErr := svc.sendbirdClient.DeleteGroupChannel(req.ChannelID)
	if deleteChannelErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errDeletingPaidGroup,
			fmt.Sprintf("Error deleting Sendbird channel for paid group chat %s and user %d: %v", req.ChannelID, currUser.ID, deleteChannelErr),
		)
	}

	commit = true

	return nil
}

func (svc *usecase) sendTaskBatch(ctx context.Context, userIDsBatch []int, channelID string, goatUserID int) error {
	message := cloudTasks.StripePaidGroupUnsubscribeTask{
		UserIDs:    userIDsBatch,
		ChannelID:  channelID,
		GoatUserID: goatUserID,
	}

	createTaskParams := cloudTasks.CreateTaskParams{
		QueueID:   cloudTasks.StripePaidGroupChatUnsubscribeQueueID,
		TargetURL: fmt.Sprintf("%s/cloudtask/v1/subscriptions/groups/batch/unsubscribe", appconfig.Config.Gcloud.APIBaseURL),
	}

	_, createTaskErr := svc.cloudTasksClient.CreateTask(ctx, createTaskParams, message)
	if createTaskErr != nil {
		return createTaskErr
	}

	return nil
}
