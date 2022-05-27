package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	userrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errCreatingFreeGroup = "Something went wrong when trying to create free group chat."

func (svc *usecase) CreateFreeGroupChannel(ctx context.Context, req request.CreateFreeGroupChannel) (*response.FreeGroup, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingFreeGroup,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errCreatingFreeGroup,
			fmt.Sprintf("CreateFreeGroupChannel returned nil user for request: %v", req),
		)
	}

	metadata := response.GroupMetadata{
		Description: req.Description,
		Benefit1:    req.Benefit1,
		Benefit2:    req.Benefit2,
		Benefit3:    req.Benefit3,
	}
	metadataBytes, marshalErr := json.Marshal(metadata)
	if marshalErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingFreeGroup,
			fmt.Sprintf("Error marshaling free group chat info during create: %v", marshalErr),
		)
	}

	if req.Name == "" {
		req.Name = user.Username
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingFreeGroup,
			fmt.Sprintf("Error creating transaction: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	if req.LinkSuffix == "" {
		newLink, generateLinkErr := svc.repo.GenerateUniqueLink(ctx, req.Name)
		if generateLinkErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				errCreatingFreeGroup,
				fmt.Sprintf("Error generating group link: %v", generateLinkErr),
			)
		}

		if newLink == nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				errCreatingFreeGroup,
				"Could not generate unique group link",
			)
		}
		req.LinkSuffix = *newLink
	}

	taken, isTakenErr := userrepo.IsLinkSuffixTaken(ctx, tx, req.LinkSuffix)
	if isTakenErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingFreeGroup,
			fmt.Sprintf("Error checking if link is taken: %v", isTakenErr),
		)
	}

	if taken.IsTaken {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"This link is unavailable.",
			"This link is unavailable.",
		)
	}

	channelData := response.GroupChannelData{
		LinkSuffix: &req.LinkSuffix,
	}

	channelDatab, marshalErr := json.Marshal(channelData)
	if marshalErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingFreeGroup,
			fmt.Sprintf("Error marshaling channel data in CreateFreeGroupChannel: %v", marshalErr),
		)
	}

	userIDs := append(req.Members, user.ID)

	createGroupChatParams := sendbird.CreateGroupChannelParams{
		UserIDs:     userIDs,
		CustomType:  constants.ChannelTypeFreeGroup,
		Name:        req.Name,
		CoverFile:   req.CoverFile,
		IsDistinct:  false, // Sendbird disallows is_distinct w/ super groups
		IsPublic:    true,
		IsSuper:     true,
		OperatorIDs: []int{user.ID},
		InviterID:   fmt.Sprint(user.ID),
		Data:        string(channelDatab),
	}

	newGroupChannel, createChannelErr := svc.sendbirdClient.CreateGroupChannel(ctx, &createGroupChatParams)
	if createChannelErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingFreeGroup,
			fmt.Sprintf("Error creating free group chat for user %d: %v", user.ID, createChannelErr),
		)
	}

	defer func() {
		if !commit {
			deleteErr := svc.sendbirdClient.DeleteGroupChannel(newGroupChannel.ChannelURL)
			if deleteErr != nil {
				vlog.Errorf(ctx, "Error in deferred delete sendbird channel %s: %v", newGroupChannel.ChannelURL, deleteErr)
			}
		}
	}()

	insertedID, insertChatErr := svc.repo.UpsertFreeGroupChat(
		ctx,
		user.ID,
		newGroupChannel.ChannelURL,
		req.LinkSuffix,
		*req.MemberLimit,
		*req.IsMemberLimitEnabled,
		metadataBytes,
		tx,
	)
	if insertChatErr != nil || insertedID == nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errCreatingFreeGroup,
			fmt.Sprintf("Error inserting free group chat for user %d: %v", user.ID, insertChatErr),
		)
	}

	res := &response.FreeGroup{
		ID:                   *insertedID,
		CreatedByUserID:      user.ID,
		ChannelID:            newGroupChannel.ChannelURL,
		Channel:              newGroupChannel,
		LinkSuffix:           req.LinkSuffix,
		IsMember:             true,
		Metadata:             metadata,
		IsMemberLimitEnabled: *req.IsMemberLimitEnabled,
		MemberLimit:          *req.MemberLimit,
	}
	commit = true
	return res, nil
}
