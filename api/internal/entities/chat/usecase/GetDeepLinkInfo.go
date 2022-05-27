package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	userrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errSorryUserDoesNotExist = "Sorry, this user doesn't seem to exist." // copying telegram's message

func (svc *usecase) GetDeepLinkInfo(ctx context.Context, req request.GetDeepLinkInfo) (*response.DeepLinkInfo, error) {
	taken, isTakenErr := userrepo.IsLinkSuffixTaken(ctx, svc.repo.MasterNode(), req.LinkSuffix)
	if isTakenErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("Error checking if link is valid: %v", isTakenErr),
		)
	}

	if !taken.IsTaken {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errSorryUserDoesNotExist,
			"Link suffix not found.",
		)
	}

	res := &response.DeepLinkInfo{}

	if taken.PaidGroupChannelID != nil {
		paidGroup, getPaidGroupHTTPErr := svc.GetPaidGroup(ctx, *taken.PaidGroupChannelID)
		if getPaidGroupHTTPErr != nil {
			return nil, getPaidGroupHTTPErr
		}
		res.PaidGroup = paidGroup.Group
		return res, nil
	} else if taken.FreeGroupChannelID != nil {
		freeGroup, getFreeGroupHTTPErr := svc.GetFreeGroup(ctx, *taken.FreeGroupChannelID)
		if getFreeGroupHTTPErr != nil {
			return nil, getFreeGroupHTTPErr
		}
		res.FreeGroup = freeGroup.Group
		return res, nil
	} else if taken.UserID != nil {
		getChannelWithUserReq := request.GetChannelWithUser{UserID: *taken.UserID}
		channel, getChannelWithUserErr := svc.GetChannelWithUser(ctx, getChannelWithUserReq)
		if getChannelWithUserErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
				fmt.Sprintf("Error getting channel with user: %v", getChannelWithUserErr),
			)
		}
		if channel != nil {
			res.UserChannelID = &channel.ChannelURL
			return res, nil
		}

		// User is valid but no channel exists, create one.
		user, getCurrentUserErr := svc.user.GetCurrentUser(ctx)
		if getCurrentUserErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
				fmt.Sprintf("Could not find user in the current context. Err: %v", getCurrentUserErr),
			)
		} else if user == nil {
			return nil, httperr.NewCtx(
				ctx,
				404,
				http.StatusNotFound,
				constants.ErrSomethingWentWrong,
				"user was nil in GetDeepLinkInfo",
			)
		}

		userIDs := []int{user.ID, *taken.UserID}
		createChannelParams := &sendbird.CreateGroupChannelParams{
			UserIDs:     userIDs,
			OperatorIDs: userIDs,
			IsDistinct:  true,
		}
		newChannel, createChannelErr := svc.sendbirdClient.CreateGroupChannel(ctx, createChannelParams)
		if createChannelErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
				fmt.Sprintf("Error creating new channel with user: %v", createChannelErr),
			)
		}
		res.UserChannelID = &newChannel.ChannelURL
		return res, nil
	}

	return nil, httperr.NewCtx(
		ctx,
		404,
		http.StatusNotFound,
		errSorryUserDoesNotExist,
		"Link suffix not found.",
	)
}
