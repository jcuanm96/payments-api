package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

func decodeStripeHandleEvent(c *fiber.Ctx) (interface{}, error) {
	signingSecret := appconfig.Config.Stripe.EventSecret
	event, eventErr := webhook.ConstructEvent([]byte(c.Body()), c.Get("Stripe-Signature"), signingSecret)
	if eventErr != nil {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Could not construct event",
			fmt.Sprintf("Error creating event: %v", eventErr),
		)
	}

	if !isValidStripeEvent(event) {
		return nil, httperr.New(
			202,
			http.StatusAccepted,
			"Stripe event received, but it was not a valid event we process",
		)
	}

	return event, nil
}

func isValidStripeEvent(event stripe.Event) bool {
	return event.Type == "charge.expired" ||
		event.Type == "charge.captured" ||
		event.Type == "invoice.paid" ||
		event.Type == "invoice.payment_failed"
}

func decodeSendbirdHandleEvent(c *fiber.Ctx) (interface{}, error) {
	token := appconfig.Config.Sendbird.MasterAPIKey

	hasher := hmac.New(sha256.New, []byte(token))
	hasher.Write(c.Body())
	calculatedHash := hex.EncodeToString(hasher.Sum(nil))

	signature := c.Get("x-sendbird-signature", "")
	if calculatedHash != signature {
		return nil, httperr.New(
			401,
			http.StatusUnauthorized,
			"Error calculating/confirming signature",
			"Error calculating/confirming signature",
		)
	}

	event := sendbird.Event{}
	unmarshalErr := json.Unmarshal(c.Body(), &event)
	if unmarshalErr != nil {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Error unmarshaling request body",
			fmt.Sprintf("Error unmarshaling request body into category event struct: %s", unmarshalErr.Error()),
		)
	}

	var specificEvent interface{}
	switch event.Category {
	case sendbird.GroupChannelCreateCategory:
		specificEvent = sendbird.GroupChannelCreateEvent{}
	case sendbird.GroupChannelChangedCategory:
		specificEvent = sendbird.GroupChannelChangedEvent{}
	case sendbird.GroupChannelJoinCategory:
		specificEvent = sendbird.GroupChannelJoinEvent{}
	case sendbird.GroupChannelInviteCategory:
		specificEvent = sendbird.GroupChannelInviteEvent{}
	default:
		return nil, httperr.New(
			202,
			http.StatusAccepted,
			"Not a category we process",
			fmt.Sprintf("Not a category we process: %s", event.Category),
		)
	}

	unmarshalSpecificErr := json.Unmarshal(c.Body(), &specificEvent)
	if unmarshalSpecificErr != nil {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Error unmarshaling request body",
			fmt.Sprintf("Error unmarshaling request body into specific event struct: %s", unmarshalSpecificErr.Error()),
		)
	}
	return specificEvent, nil
}
