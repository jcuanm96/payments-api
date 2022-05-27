package service

import (
	"context"
	"errors"
	"fmt"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) SendGoatPostNotification(ctx context.Context, goatID int, postID int) error {
	goatUser, getGoatUserErr := svc.user.GetUserByID(ctx, goatID)
	if getGoatUserErr != nil {
		vlog.Errorf(ctx, "Could not get creator user: %v", getGoatUserErr)
		return getGoatUserErr
	} else if goatUser == nil {
		return errors.New("creator user was nil")
	}

	title := "Vama"
	body := fmt.Sprintf("%s %s just shared a post", goatUser.FirstName, goatUser.LastName)
	notificationErr := svc.fcmClient.GoatPostPushNotification(ctx, postID, goatID, title, body)
	return notificationErr
}
