package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	userrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errCreatingChannel = "Something went wrong creating channel."

/* CreateGoatChatChannel checks if a channel already exists with the other user. If so, return it.
If the chat is not between a USER and a CREATOR, chat type is NORMAL.
Otherwise, check if the USER is a contact of the CREATOR.
If the USER is a contact, chat type is NORMAL, else it's CREATOR.
Then create the channel with the appropriate channel data.
*/
func (svc *usecase) CreateGoatChatChannel(ctx context.Context, req request.CreateGoatChatChannel) (*response.CreateGoatChatChannel, error) {
	existingChannel, checkExistingChannelErr := svc.GetChannelWithUser(ctx, request.GetChannelWithUser{UserID: req.OtherUserID})
	if checkExistingChannelErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingChannel,
			fmt.Sprintf("Error checking if channel exists with %d: %v", req.OtherUserID, checkExistingChannelErr),
		)
	}

	if existingChannel != nil {
		res := &response.CreateGoatChatChannel{
			Channel: existingChannel,
		}
		return res, nil
	}

	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingChannel,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errCreatingChannel,
			fmt.Sprintf("User was nil. Request: %v", req),
		)
	}

	otherUser, otherUserErr := svc.user.GetUserByID(ctx, req.OtherUserID)
	if otherUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingChannel,
			fmt.Sprintf("Could not find user in the current context. Err: %v", otherUserErr),
		)
	}
	if otherUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errCreatingChannel,
			fmt.Sprintf("User %d was nil", req.OtherUserID),
		)
	}

	category, getCategoryErr := svc.getChatCategory(ctx, user, otherUser)
	if getCategoryErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingChannel,
			fmt.Sprintf("Error getting chat category between %d and %d: %v", user.ID, otherUser.ID, getCategoryErr),
		)
	}

	// if category is GOAT then set initial creator chat data
	channelData := "{}"
	if category == constants.ChannelTypePaidDirect {
		state := constants.ChannelStateActive
		data := response.ChannelData{
			ChannelState: &state,
		}
		dataBytes, marshalErr := json.Marshal(data)
		if marshalErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				errCreatingChannel,
				fmt.Sprintf("Error marshaling channel data: %v", marshalErr),
			)
		}

		channelData = string(dataBytes)
	}

	userIDs := []int{user.ID, otherUser.ID}
	createGroupChannelParams := sendbird.CreateGroupChannelParams{
		UserIDs:     userIDs,
		CustomType:  category,
		Data:        channelData,
		IsDistinct:  true,
		OperatorIDs: userIDs,
	}
	channel, createChannelErr := svc.sendbirdClient.CreateGroupChannel(ctx, &createGroupChannelParams)
	if createChannelErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingChannel,
			fmt.Sprintf("Error creating group channel: %v", createChannelErr),
		)
	}

	res := &response.CreateGoatChatChannel{
		Channel: channel,
	}

	return res, nil
}

func (svc *usecase) getChatCategory(ctx context.Context, user, otherUser *response.User) (string, error) {
	var userID int
	var goatID int
	if user.Type == "USER" && otherUser.Type == "GOAT" {
		userID = user.ID
		goatID = otherUser.ID
	} else if user.Type == "GOAT" && otherUser.Type == "USER" {
		userID = otherUser.ID
		goatID = user.ID
	} else {
		return constants.ChannelTypeDirect, nil
	}

	// do we need to check for blocked?

	isContact, isContactErr := userrepo.IsContact(ctx, svc.repo.MasterNode(), goatID, userID)
	if isContactErr != nil {
		vlog.Errorf(ctx, "Error checking if %d is contact of %d: %v", userID, goatID, isContactErr)
		return "", isContactErr
	}

	if isContact {
		return constants.ChannelTypeDirect, nil
	} else {
		return constants.ChannelTypePaidDirect, nil
	}
}
