package request

type SubscribeCurrUserToGoat struct {
	GoatUserID int `json:"goatUserID"`
}

type GetSubscriptionPrices struct {
	GoatUserID int `json:"goatUserID"`
}

type UnsubscribeCurrUserFromGoat struct {
	GoatUserID int `json:"goatUserID"`
}

type GetUserSubscriptions struct {
	Limit    *uint64 `json:"limit"`
	CursorID *int    `json:"cursorID"`
}

type UpdateSubscriptionPrice struct {
	PriceInSmallestDenom int64  `json:"priceInSmallestDenom"`
	Currency             string `json:"currency"`
	// Later add Tier here
}

type CheckUserSubscribedToGoat struct {
	ProviderUserID int `json:"providerUserID"`
}

type GetMyPaidGroupSubscriptions struct {
	Limit    int64 `json:"limit"`
	CursorID int   `json:"cursorID"`
}

type GetMyPaidGroupSubscription struct {
	ChannelID string `json:"channelID"`
}
