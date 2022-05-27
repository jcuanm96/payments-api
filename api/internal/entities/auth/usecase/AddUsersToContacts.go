package service

import (
	"context"
	"fmt"
	"time"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	cloudTasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) AddUsersToContacts(ctx context.Context, exec utils.Executable, user response.User) error {
	pendingContactsIDs, getPendingContactsErr := svc.repo.GetPendingContactsByPhone(ctx, user.Phonenumber)
	if getPendingContactsErr != nil {
		return getPendingContactsErr
	} else if len(pendingContactsIDs) == 0 {
		return nil
	}

	batchSize := 50
	contactUserIDs := []int{}
	for _, contactUserID := range pendingContactsIDs {
		contactUserIDs = append(contactUserIDs, contactUserID)
		if len(contactUserIDs) >= batchSize {
			sendTaskBatchErr := svc.sendTaskBatch(ctx, contactUserIDs, user.ID)
			if sendTaskBatchErr != nil {
				sendTaskBatchErrMsg := fmt.Sprintf("Something went wrong when sending task for AddUsersContacts by userID: %d. ContactsUserIDs: %v. Err: %s", user.ID, contactUserIDs, sendTaskBatchErr.Error())
				vlog.Errorf(ctx, sendTaskBatchErrMsg)
				telegram.TelegramClient.SendMessage(sendTaskBatchErr.Error())
			}
			contactUserIDs = []int{}
		}
	}

	if len(contactUserIDs) > 0 {
		sendTaskBatchErr := svc.sendTaskBatch(ctx, contactUserIDs, user.ID)
		if sendTaskBatchErr != nil {
			sendTaskBatchErrMsg := fmt.Sprintf("Error creating task as final element in AddUserContacts by user %d. ContactsUserIDs: %v. Err: %s", user.ID, contactUserIDs, sendTaskBatchErr.Error())
			vlog.Errorf(ctx, sendTaskBatchErrMsg)
			telegram.TelegramClient.SendMessage(sendTaskBatchErr.Error())
		}
	}

	return nil
}

func (svc *usecase) sendTaskBatch(ctx context.Context, contactUserIDs []int, userID int) error {
	message := cloudTasks.AddUserContactsTask{
		UserID:         userID,
		ContactUserIDs: contactUserIDs,
	}

	scheduleTime := time.Now().Add(time.Minute * 1)
	createTaskParams := cloudTasks.CreateTaskParams{
		QueueID:      cloudTasks.AddUserContactsQueueID,
		TargetURL:    fmt.Sprintf("%s/cloudtask/v1/contacts/me/batch", appconfig.Config.Gcloud.APIBaseURL),
		ScheduleTime: &scheduleTime,
	}

	_, createTaskErr := svc.cloudTasksClient.CreateTask(ctx, createTaskParams, message)
	if createTaskErr != nil {
		return createTaskErr
	}

	return nil
}
