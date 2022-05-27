package response

import "time"

type GetGoatChatPrice struct {
	PriceInSmallestDenom int64  `json:"priceInSmallestDenom"`
	Currency             string `json:"currency"`
}

type TransactionItem struct {
	ID        int    `json:"id"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	CreatedAt int64  `json:"createdAt"`
	User      *User  `json:"user"`
	Fees      int64  `json:"fees"`
	Type      string `json:"type"`
}

type GetTransactions struct {
	Transactions []TransactionItem `json:"transactions"`
}

type DefaultPaymentMethod struct{}

type GetBalance struct {
	Currency   string    `json:"currency"`
	Amount     int64     `json:"amount"`
	Bank       *Bank     `json:"bank"`
	NextPayout time.Time `json:"nextPayout"`
}

type ConfirmPaymentIntent struct {
	ChargeID string `json:"chargeID"`
}

type UpsertGoatChatsPrice struct{}

type PayoutToBank struct {
	ArrivalDate int64  `json:"arrivalDate"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
}

type GetMyPaymentMethods struct {
	PaymentMethods []PaymentMethod `json:"paymentMethods"`
}

type PaymentMethod struct {
	Card *Card `json:"card"`
	Bank *Bank `json:"bank"`
}

type Card struct {
	Brand    string `json:"brand"`
	ExpMonth uint8  `json:"expMonth"`
	ExpYear  uint16 `json:"expYear"`
	Funding  string `json:"funding"`
	Last4    string `json:"last4"`
	Country  string `json:"country"`
	CVCCheck string `json:"cvcCheck"`
	Name     string `json:"name"`
}

type Bank struct {
	BankName          string         `json:"bankName"`
	AccountNumber     string         `json:"accountNumber"`
	RoutingNumber     string         `json:"routingNumber,omitempty"`
	AccountType       string         `json:"accountType"`
	AccountHolderName string         `json:"accountHolderName"`
	AccountHolderType string         `json:"accountHolderType"`
	Currency          string         `json:"currency"`
	Country           string         `json:"country"`
	BillingAddress    BillingAddress `json:"billingAddress"`
}

type BillingAddress struct {
	Street1    string `json:"street1"`
	Street2    string `json:"street2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type ProviderPaymentInfo struct {
	ID                 int       `json:"id"`
	FirstName          string    `json:"firstName"`
	LastName           string    `json:"lastName"`
	BalanceOwed        int64     `json:"balanceOwed"`
	Currency           string    `json:"currency"`
	LastPayoutTS       time.Time `json:"lastPayoutTS"`
	LastPayoutPeriodID int       `json:"lastPayoutPeriodID"`
	Bank               Bank      `json:"bank"`
}

type ListUnpaidProviders struct {
	Providers []ProviderPaymentInfo `json:"providers"`
}

type BankInfo struct {
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	RoutingNumber string `json:"routingNumber"`
	Type          string `json:"type"`
}

type GetPayoutPeriods struct {
	PayoutPeriods []PayoutPeriod `json:"payoutPeriods"`
}

type PayoutPeriod struct {
	ID      int   `json:"id"`
	StartTS int64 `json:"startTS"`
	EndTS   int64 `json:"endTS"`
}

type ListPayoutHistory struct {
	PayoutHistory []PayoutHistoryDatum `json:"payoutHistory"`
	HasNext       bool                 `json:"hasNext"`
}

type PayoutHistoryDatum struct {
	PayoutPeriod PayoutPeriod `json:"payoutPeriod"`
	BalanceOwed  int64        `json:"balanceOwed"`
	BalancePaid  int64        `json:"balancePaid"`
}
