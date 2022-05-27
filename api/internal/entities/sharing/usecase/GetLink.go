package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	chatrepo "github.com/VamaSingapore/vama-api/internal/entities/chat/repositories"
	"github.com/VamaSingapore/vama-api/internal/entities/sharing"
	userrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func (svc *usecase) GetLink(ctx context.Context, req request.GetLink) (*response.GetVamaMeLink, error) {
	runnable := svc.repo.MasterNode()
	taken, isTakenErr := userrepo.IsLinkSuffixTaken(ctx, runnable, req.LinkSuffix)
	if isTakenErr != nil {
		vlog.Errorf(ctx, "Error checking if link suffix %s is valid: %v", req.LinkSuffix, isTakenErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	} else if !taken.IsTaken {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"We couldn't find anything for this link.",
		)
	}

	deepLinkBaseURL := appconfig.Config.Gcloud.RedirectBaseURL
	deepLink := fmt.Sprintf("https://%s/%s", deepLinkBaseURL, req.LinkSuffix)
	dynamicLink, getDynamicLinkErr := sharing.CreateDynamicLink(ctx, deepLink)
	if getDynamicLinkErr != nil {
		vlog.Errorf(ctx, "Error creating dynamic link for suffix %s: %v", req.LinkSuffix, getDynamicLinkErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	}

	res := &response.GetVamaMeLink{
		DynamicLink: dynamicLink,
	}

	if taken.UserID != nil {
		// When biolink code is set up, use this.
		// bioLinks, getBioLinksErr := svc.repo.GetBioLinks(ctx, *taken.UserID)
		// if getBioLinksErr != nil {
		// 	return nil, httperr.NewCtx(
		// 	    ctx,
		// 		500,
		// 		http.StatusInternalServerError,
		// 		constants.ErrSomethingWentWrong,
		// 	)
		// }
		// res.BioLinks = bioLinks
		bioData, getBioDataErr := svc.repo.GetBioData(ctx, runnable, strings.ToLower(req.LinkSuffix), *taken.UserID)
		if getBioDataErr != nil {
			vlog.Errorf(ctx, "Error getting bio data for user %s: %v", req.LinkSuffix, getBioDataErr)
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
			)
		}
		res.BioData = bioData
	} else if taken.PaidGroupChannelID != nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			"Chat not found.",
		)
		// groupInfo, getGroupInfoErr := chatrepo.GetPaidGroup(ctx, svc.repo.MasterNode(), *taken.PaidGroupChannelID)
		// if getGroupInfoErr != nil {
		// 	vlog.Errorf(ctx, "Error getting group info for link %s: %v", req.LinkSuffix, getGroupInfoErr)
		// 	return nil, httperr.NewCtx(
		// 		ctx,
		// 		500,
		// 		http.StatusInternalServerError,
		// 		constants.ErrSomethingWentWrong,
		// 	)
		// }

		// channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(groupInfo.ChannelID, sendbird.GetGroupChannelParams{})
		// if getChannelErr != nil {
		// 	vlog.Errorf(ctx, "Error getting sendbird channel %s in GetLink: %v", groupInfo.ChannelID, getChannelErr)
		// 	return nil, httperr.NewCtx(
		// 		ctx,
		// 		500,
		// 		http.StatusInternalServerError,
		// 		constants.ErrSomethingWentWrong,
		// 	)
		// }

		// goatUser, getGoatUserErr := svc.user.GetUserByID(ctx, groupInfo.GoatID)
		// if getGoatUserErr != nil {
		// 	vlog.Errorf(ctx, "Error getting user %d by ID: %v", groupInfo.GoatID, getGoatUserErr)
		// 	return nil, httperr.NewCtx(
		// 		ctx,
		// 		500,
		// 		http.StatusInternalServerError,
		// 		constants.ErrSomethingWentWrong,
		// 	)
		// }

		// // Only keep the bare minimum fields
		// returnedGoatUser := &response.User{
		// 	Username:      goatUser.Username,
		// 	ProfileAvatar: goatUser.ProfileAvatar,
		// 	FirstName:     goatUser.FirstName,
		// 	LastName:      goatUser.LastName,
		// }

		// redirectPaidGroupInfo := &response.RedirectPaidGroupInfo{
		// 	Name:                 channel.Name,
		// 	PriceInSmallestDenom: groupInfo.PriceInSmallestDenom,
		// 	Currency:             groupInfo.Currency,
		// 	Metadata:             groupInfo.Metadata,
		// 	MemberCount:          channel.JoinedMemberCount,
		// 	ProfileAvatar:        channel.CoverURL,
		// 	Goat:                 returnedGoatUser,
		// }
		// res.PaidGroupInfo = redirectPaidGroupInfo
	} else if taken.FreeGroupChannelID != nil {
		groupInfo, getGroupInfoErr := chatrepo.GetFreeGroup(ctx, svc.repo.MasterNode(), *taken.FreeGroupChannelID)
		if getGroupInfoErr != nil {
			vlog.Errorf(ctx, "Error getting free group info for link %s: %v", req.LinkSuffix, getGroupInfoErr)
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
			)
		}

		channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(groupInfo.ChannelID, sendbird.GetGroupChannelParams{})
		if getChannelErr != nil {
			vlog.Errorf(ctx, "Error getting sendbird channel %s in GetLink: %v", groupInfo.ChannelID, getChannelErr)
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
			)
		}

		createdByUser, getCreatedByUserErr := svc.user.GetUserByID(ctx, groupInfo.CreatedByUserID)
		if getCreatedByUserErr != nil {
			vlog.Errorf(ctx, "Error getting user %d by ID: %v", groupInfo.CreatedByUserID, getCreatedByUserErr)
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				constants.ErrSomethingWentWrong,
			)
		}

		// Only keep the bare minimum fields
		returnedCreatedByUser := &response.User{
			Username:      createdByUser.Username,
			ProfileAvatar: createdByUser.ProfileAvatar,
			FirstName:     createdByUser.FirstName,
			LastName:      createdByUser.LastName,
		}

		redirectFreeGroupInfo := &response.RedirectFreeGroupInfo{
			Name:          channel.Name,
			Metadata:      groupInfo.Metadata,
			MemberCount:   channel.JoinedMemberCount,
			ProfileAvatar: channel.CoverURL,
			Goat:          returnedCreatedByUser,
		}
		res.FreeGroupInfo = redirectFreeGroupInfo
	} else {
		vlog.Errorf(ctx, "Link %s is taken but userID and channelID both nil", req.LinkSuffix)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	}

	return res, nil
}
