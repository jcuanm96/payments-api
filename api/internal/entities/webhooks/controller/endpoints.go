package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/webhooks"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	telegram "github.com/VamaSingapore/vama-api/internal/telegramClient"
	"github.com/stripe/stripe-go/v72"
)

func StripeHandleEvent(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(webhooks.Usecase)
	stripeEvent := incomeRequest.(stripe.Event)

	handlerErr := svc.StripeHandleEvent(ctx, stripeEvent)
	if handlerErr != nil {
		handlerErrMsg := fmt.Sprintf("StripeHandleEvent returned an error for Stripe id %s and ts %d. Err: %s", stripeEvent.ID, stripeEvent.Created, handlerErr.Error())
		telegram.TelegramClient.SendMessage(handlerErrMsg)
	}

	return nil, handlerErr
}

func SendbirdHandleEvent(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(webhooks.Usecase)

	reqBytes, marshalErr := json.Marshal(incomeRequest)
	if marshalErr != nil {
		return nil, marshalErr
	}

	categoryEvent := sendbird.Event{}
	unmarshalErr := json.Unmarshal(reqBytes, &categoryEvent)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	var handlerErr error
	switch categoryEvent.Category {
	case sendbird.GroupChannelCreateCategory:
		event := &sendbird.GroupChannelCreateEvent{}
		unmarshalErr := json.Unmarshal(reqBytes, event)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		handlerErr = svc.SendbirdCreateGroupEvent(ctx, event)

	case sendbird.GroupChannelChangedCategory:
		event := &sendbird.GroupChannelChangedEvent{}
		unmarshalErr := json.Unmarshal(reqBytes, event)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		handlerErr = svc.SendbirdChangedGroupEvent(ctx, event)
	case sendbird.GroupChannelJoinCategory:
		event := &sendbird.GroupChannelJoinEvent{}
		unmarshalErr := json.Unmarshal(reqBytes, event)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		handlerErr = svc.SendbirdJoinGroupEvent(ctx, event)
	case sendbird.GroupChannelInviteCategory:
		event := &sendbird.GroupChannelInviteEvent{}
		unmarshalErr := json.Unmarshal(reqBytes, event)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		handlerErr = svc.SendbirdInviteGroupEvent(ctx, event)
	}

	if handlerErr != nil {
		handlerErrMsg := fmt.Sprintf("HandleEvent returned an error for sendbird webhook event. Err: %s", handlerErr.Error())
		telegram.TelegramClient.SendMessage(handlerErrMsg)
		return nil, handlerErr
	}

	return nil, nil
}
