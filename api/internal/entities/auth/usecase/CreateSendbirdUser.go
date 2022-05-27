package service

import (
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

func (svc *usecase) CreateSendBirdUser(currUser response.User) (string, error) {
	createUserParams := &sendbird.CreateUserParams{
		UserID:           currUser.ID,
		ProfileURL:       "",
		Nickname:         fmt.Sprintf("%s %s", currUser.FirstName, currUser.LastName),
		IssueAccessToken: true,
	}

	sendbirdUser, createUserErr := svc.sendbirdClient.CreateUser(createUserParams)
	if createUserErr != nil {
		return "", createUserErr
	}

	if sendbirdUser.AccessToken == "" {
		return "", fmt.Errorf("SendBird access token came back nil for user %d", currUser.ID)
	}

	return sendbirdUser.AccessToken, nil
}
