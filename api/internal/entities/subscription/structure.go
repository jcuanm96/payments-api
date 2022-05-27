package subscription

import (
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

func FormatSubscriptionProduct(username string, tierName string, price int64, currency string) string {
	return fmt.Sprintf("@%s-%s-%d%s", username, tierName, price, currency)
}

type PaidGroupChatInfo struct {
	ChannelID            string
	GoatID               int
	StripeProductID      string
	PriceInSmallestDenom int64
	Currency             string
	LinkSuffix           string
	Metadata             response.GroupMetadata
	MemberLimit          int
	IsMemberLimitEnabled bool
}
