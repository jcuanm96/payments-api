package vamabot

import (
	"context"
	"fmt"
	"strings"

	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (c *client) SendWelcomeMessages(ctx context.Context, userID int) error {
	ids := []int{c.vamaUserID, userID}
	channel, createChannelErr := c.sendbirdClient.CreateGroupChannel(ctx, &sendbird.CreateGroupChannelParams{
		UserIDs:     ids,
		OperatorIDs: ids,
		IsDistinct:  true,
	})
	if createChannelErr != nil {
		vlog.Errorf(ctx, "Error creating welcome channel for userID %d: %s", userID, createChannelErr.Error())
		return createChannelErr
	}

	sendMessageParams := &sendbird.SendMessageParams{
		UserID:      fmt.Sprint(c.vamaUserID),
		Message:     strings.TrimSpace(welcomeMessage),
		MessageType: "MESG",
	}
	_, messageErr := c.sendbirdClient.SendMessage(channel.ChannelURL, sendMessageParams)
	if messageErr != nil {
		vlog.Errorf(ctx, "Error sending message to userID %d in channel %s: %s", userID, channel.ChannelURL, messageErr.Error())
		return messageErr
	}

	return nil
}

const welcomeMessage = `
Welcome to Vama! 👾
👉 Here’s some help to get started!

We invented a brand new way to message your favorite creators and receive a guaranteed response.

Get started for free by adding your friends, creating groups, and using our insanely fast chat service. 💬  

Vama is a community for Crypto, NFTs, Memes, (or anything) and privacy is our top priority. 🤫

When you message a creator on Vama, they will respond directly (or your money back!) 💎🙌

Everyone’s an expert.

Made with love, Vama Team ☕
`
