package wallet

import (
	"context"
	"fmt"
	"strconv"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

type LedgerEntry struct {
	ProviderUserID      int
	CustomerUserID      int
	StripeTransactionID string
	SourceType          string
	Amount              int64
	StripeFee           int64
	VamaFee             int64
	CreatedTS           int64
	BalanceDelta        int64
	Currency            string
	PayPeriodID         int
}

// ID can be either:
// ch_...
// in_...
type StripeTxnEvent struct {
	ID        string
	CreatedAt int64
	Metadata  *StripeEventMetadata
}

type StripeEventMetadata struct {
	ProviderUserID int
	CustomerUserID int
	IsTrial        *bool
	TierID         *int
	ChannelID      *string
}

type PendingBalanceNotification struct {
	FcmToken         string
	AvailableBalance int
}

func GetUserIDsFromStripeMetadata(ctx context.Context, metadata map[string]string) (*StripeEventMetadata, error) {
	providerUserID, providerUserIDErr := GetIDFromMetadata(ctx, metadata, "providerUserID")
	if providerUserIDErr != nil {
		return nil, providerUserIDErr
	}
	customerUserID, customerUserIDErr := GetIDFromMetadata(ctx, metadata, "customerUserID")
	if customerUserIDErr != nil {
		return nil, customerUserIDErr
	}
	res := &StripeEventMetadata{
		ProviderUserID: providerUserID,
		CustomerUserID: customerUserID,
	}
	return res, nil
}

// key is either one of the string literals "providerUserID" or "customerUserID"
func GetIDFromMetadata(ctx context.Context, metadata map[string]string, key string) (int, error) {
	userIDStr, ok := metadata[key]
	if !ok {
		return -1, fmt.Errorf("no %s in metadata", key)
	}
	userID, atoiErr := strconv.Atoi(userIDStr)
	if atoiErr != nil {
		vlog.Errorf(ctx, "Error converting %sStr %s to int: %s", key, userIDStr, atoiErr.Error())
		return -1, atoiErr
	}
	return userID, nil
}
