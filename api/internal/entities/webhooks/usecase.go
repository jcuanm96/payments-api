package webhooks

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/wallet"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/stripe/stripe-go/v72"
)

type Usecase interface {
	// Stripe
	StripeHandleEvent(ctx context.Context, stripeEvent stripe.Event) error
	StripePendingChargeEvent(ctx context.Context, stripeEvent stripe.Event) error
	StripeSubscriptionPaymentSuccess(ctx context.Context, event stripe.Event) (*wallet.StripeEventMetadata, error)
	StripeSubscriptionPaymentFailure(ctx context.Context, event stripe.Event) (*response.HandleSubscriptionPaymentEvent, error)

	// Sendbird
	SendbirdCreateGroupEvent(ctx context.Context, event *sendbird.GroupChannelCreateEvent) error
	SendbirdChangedGroupEvent(ctx context.Context, event *sendbird.GroupChannelChangedEvent) error
	SendbirdJoinGroupEvent(ctx context.Context, event *sendbird.GroupChannelJoinEvent) error
	SendbirdInviteGroupEvent(ctx context.Context, event *sendbird.GroupChannelInviteEvent) error
}
