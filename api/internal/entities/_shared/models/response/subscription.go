package response

import (
	"time"
)

type GoatSubscriptionInfo struct {
	ID                   int    `json:"id"`
	GoatUserID           int    `json:"goatUserID"`
	PriceInSmallestDenom int64  `json:"priceInSmallestDenom"`
	Currency             string `json:"currency"`
	TierName             string `json:"tierName"`
	StripeProductID      string `json:"stripeProductID"`
}

type SubscribeCurrUserToGoat struct{}

type UpdateSubscriptionPrice struct{}

type UserSubscription struct {
	ID                   int       `json:"id"`
	CurrentPeriodEnd     time.Time `json:"currentPeriodEnd"`
	UserID               int       `json:"userID"`
	StripeSubscriptionID string    `json:"stripeSubscriptionID"`
	GoatUser             User      `json:"goatUser"`
	TierID               int       `json:"tierID"`
	IsRenewing           bool      `json:"isRenewing"`
}

type GetUserSubscriptions struct {
	LastID        int                `json:"lastID"`
	Subscriptions []UserSubscription `json:"subscriptions"`
}

type UnsubscribeCurrUserFromGoat struct{}

type CheckUserSubscribedToGoat struct {
	IsSubscribed     bool      `json:"isSubscribed"`
	IsRenewing       bool      `json:"isRenewing"`
	CurrentPeriodEnd time.Time `json:"currentPeriodEnd"`
}

type PaidGroupChatSubscription struct {
	ID                   int        `json:"id"`
	CurrentPeriodEnd     time.Time  `json:"currentPeriodEnd"`
	UserID               int        `json:"userID,omitempty"`
	StripeSubscriptionID string     `json:"stripeSubscriptionID,omitempty"`
	GoatUser             User       `json:"goatUser"`
	ChannelID            string     `json:"channelID"`
	Group                *PaidGroup `json:"group"`
	IsRenewing           bool       `json:"isRenewing"`
}

type GetMyPaidGroupSubscriptions struct {
	Subscriptions []PaidGroupChatSubscription `json:"subscriptions"`
}

type GetMyPaidGroupSubscription struct {
	Subscription *PaidGroupChatSubscription `json:"subscription"`
}
