package controller

import (
	"context"
	"fmt"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/contact"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
)

func CreateContact(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(contact.Usecase)
	req := incomeRequest.(request.CreateContact)

	res, err := svc.CreateContact(ctx, req.ContactID)

	return res, err
}

func DeleteContact(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(contact.Usecase)
	req := incomeRequest.(request.DeleteContact)

	res, err := svc.DeleteContact(ctx, req.ContactID)

	return res, err
}

func GetContacts(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(contact.Usecase)
	req := incomeRequest.(request.GetContacts)

	res, err := svc.GetContacts(ctx, req)

	return res, err
}

func IsContact(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(contact.Usecase)
	req := incomeRequest.(request.IsContact)

	res, err := svc.IsContact(ctx, req.ContactID)

	return res, err
}

func UploadContacts(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(contact.Usecase)
	req := incomeRequest.(request.UploadContacts)

	err := svc.UploadContacts(ctx, req)
	return nil, err
}

func GetRecommendations(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(contact.Usecase)

	limit := 50
	res, err := svc.GetRecommendations(ctx, limit)
	return res, err
}

func BatchAddUsersToContacts(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(contact.Usecase)

	req := incomeRequest.(cloudtasks.AddUserContactsTask)

	scheduledErr := svc.BatchAddUsersToContacts(ctx, req)
	if scheduledErr != nil {
		scheduledErrMsg := fmt.Sprintf("Error occurred when processing task from BatchAddUsersToContacts for user %d. ContactUserIDs: %v. Err: %s", req.UserID, req.ContactUserIDs, scheduledErr.Error())
		telegram.TelegramClient.SendMessage(scheduledErrMsg)
	}

	return nil, scheduledErr
}
