package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	subscriptionrepo "github.com/VamaSingapore/vama-api/internal/entities/subscription/repositories"
	"github.com/VamaSingapore/vama-api/internal/entities/webhooks"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const errBanningUser = "Something went wrong when trying to ban user. Please try again."

func (svc *usecase) BanUserFromPaidGroup(ctx context.Context, bannedUserID int, channelID string) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errBanningUser,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errBanningUser,
			"user was nil in BanUserFromPaidGroup",
		)
	}

	getGroupChannelParams := sendbird.GetGroupChannelParams{
		ShowMember: true,
	}
	channel, getChannelErr := svc.sendbirdClient.GetGroupChannel(channelID, getGroupChannelParams)
	if getChannelErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errBanningUser,
			fmt.Sprintf("Error getting Sendbird channel for banned user %d and channel %s: %v", bannedUserID, channelID, getChannelErr),
		)
	} else if channel == nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errBanningUser,
			fmt.Sprintf("Error the Sendbird channel %s came back as nil.", channelID),
		)
	}

	if !channel.HasOperator(fmt.Sprint(user.ID)) {
		return httperr.NewCtx(
			ctx,
			403,
			http.StatusForbidden,
			"You cannot ban a user in this channel.",
			"This user cannot ban a user in this channel.",
		)
	}

	tx, txErr := svc.repo.MasterNode().Begin(ctx)
	if txErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errBanningUser,
			fmt.Sprintf("Error starting transaction: %v", txErr),
		)
	}

	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	var wg sync.WaitGroup

	var banUserErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		banUserParams := sendbird.BanUserFromGroupChannelParams{
			UserID:  fmt.Sprint(bannedUserID),
			AgentID: fmt.Sprint(user.ID),
		}
		banUserErr = svc.sendbirdClient.BanUserFromGroupChannel(channelID, &banUserParams)
	}()

	var subscription *response.PaidGroupChatSubscription
	var getSubErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		subscription, getSubErr = subscriptionrepo.GetPaidGroupSubscription(ctx, tx, bannedUserID, channelID)
	}()

	var bannedUser *response.User
	var getBannedUserErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		bannedUser, getBannedUserErr = svc.user.GetUserByID(ctx, bannedUserID)
	}()

	wg.Wait()

	if banUserErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errBanningUser,
			fmt.Sprintf("Error banning user %d from channel %s: %v", bannedUserID, channelID, banUserErr),
		)
	} else if getSubErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errBanningUser,
			fmt.Sprintf("Error getting subscription for banned user %d and channel %s: %v", bannedUserID, channelID, getSubErr),
		)
	}

	if subscription != nil {
		_, cancelStripeSubErr := svc.stripeClient.Subscriptions.Cancel(subscription.StripeSubscriptionID, nil)
		if cancelStripeSubErr != nil {
			return httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				errBanningUser,
				fmt.Sprintf("Error cancelling subscription for banned user %d and channel %s: %v", bannedUserID, channelID, cancelStripeSubErr),
			)
		}

		deleteSubErr := svc.subscription.DeletePaidGroupSubscription(ctx, tx, channelID, bannedUserID)
		if deleteSubErr != nil {
			return httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				errBanningUser,
				fmt.Sprintf("Error deleting subscription for banned user %d and channel %s. Err: %v", bannedUserID, channelID, deleteSubErr),
			)
		}
	}

	if getBannedUserErr != nil {
		vlog.Errorf(ctx, "Error getting the banned user %d in BanUserFromPaidGroup with channel %s: %v", bannedUserID, channelID, getBannedUserErr)
		bannedUser = &response.User{}
	} else if bannedUser == nil {
		vlog.Errorf(ctx, "The banned user %d in BanUserFromPaidGroup returned nil in channel %s.", bannedUserID, channelID)
		bannedUser = &response.User{}
	}

	insertErr := svc.repo.InsertBannedChatUser(ctx, bannedUserID, user.ID, channelID, tx)
	if insertErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errBanningUser,
			fmt.Sprintf("Error inserting banned user %d into table for channel %s by user %d: %v", bannedUserID, channelID, user.ID, insertErr),
		)
	}

	// Send message to channel
	message, data, rangesErr := webhooks.CalculateRemovedGroupAdminMessageRanges(user, bannedUser)
	if rangesErr != nil {
		vlog.Errorf(ctx, "Error calculating ranges for banned group message: %v", rangesErr)
		return nil
	}

	sendMessageErr := svc.SendGroupMessage(ctx, channelID, user.ID, message, data, constants.ChannelEventRemoveCustomType)
	if sendMessageErr != nil {
		vlog.Errorf(ctx, "Error sending banned group message in %s: %v", channelID, sendMessageErr)
	}

	pushErr := svc.push.SendRemovedGroupNotification(ctx, bannedUser, user, channel)
	if pushErr != nil {
		vlog.Errorf(ctx, "Error sending banned group push to %d: %v", bannedUser.ID, pushErr)
	}

	commit = true
	return nil
}
