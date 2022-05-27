package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	cloudTasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	subscriptionrepo "github.com/VamaSingapore/vama-api/internal/entities/subscription/repositories"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errCancellingPaidGroup = "Something went wrong when trying to ban user. Please try again."

func (svc *usecase) CancelPaidGroup(ctx context.Context, req request.CancelPaidGroup) (*response.CancelPaidGroup, error) {
	currUser, getCurrUserErr := svc.user.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCancellingPaidGroup,
			fmt.Sprintf("Could not find user in the current context. Err: %v", getCurrUserErr),
		)
	} else if currUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errCancellingPaidGroup,
			"user was nil in CancelPaidGroup",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCancellingPaidGroup,
			fmt.Sprintf("Error creating transaction: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	oldSubscription, getSubErr := subscriptionrepo.GetPaidGroupSubscription(ctx, tx, currUser.ID, req.ChannelID)
	if getSubErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCancellingPaidGroup,
			fmt.Sprintf("Error getting paid group chat subscription for user %d and channel %s. Err: %v", currUser.ID, req.ChannelID, getSubErr),
		)
	}

	var scheduleTime *time.Time
	if oldSubscription == nil || !oldSubscription.IsRenewing || time.Now().After(oldSubscription.CurrentPeriodEnd) {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			errCancellingPaidGroup,
			fmt.Sprintf("User %d tried to cancel group %s but they are not subscribed", currUser.ID, req.ChannelID),
		)
	} else {
		// Give buffer after the period end to account for delays in resubscriptions processing, etc.
		removeTime := oldSubscription.CurrentPeriodEnd.Add(2 * time.Hour)
		scheduleTime = &removeTime
	}

	cancelledSubscription, unsubscribeErr := svc.subscription.UnsubscribePaidGroup(ctx, tx, req.ChannelID, currUser)
	if unsubscribeErr != nil {
		return nil, unsubscribeErr
	}

	message := cloudTasks.RemoveFromPaidGroupTask{
		UserID:    currUser.ID,
		ChannelID: req.ChannelID,
	}

	// Cloud tasks will not let you schedule more than 30 days in the future.
	days := 30
	maxScheduleTime := time.Now().AddDate(0, 0, days)
	if scheduleTime.After(maxScheduleTime) {
		scheduleTime = &maxScheduleTime
	}
	createTaskParams := cloudTasks.CreateTaskParams{
		QueueID:      cloudTasks.RemoveFromPaidGroupQueueID,
		TargetURL:    fmt.Sprintf("%s/cloudtask/v1/chat/paid/group/remove", appconfig.Config.Gcloud.APIBaseURL),
		ScheduleTime: scheduleTime,
	}

	_, createTaskErr := svc.cloudTasksClient.CreateTask(ctx, createTaskParams, message)
	if createTaskErr != nil {
		cancelErrMsg := fmt.Sprintf("Error creating task to remove %d from paid group %s: %v", currUser.ID, req.ChannelID, createTaskErr)
		telegram.TelegramClient.SendMessage(cancelErrMsg)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCancellingPaidGroup,
			cancelErrMsg,
		)
	}

	commit = true
	return &response.CancelPaidGroup{Subscription: cancelledSubscription}, nil
}
