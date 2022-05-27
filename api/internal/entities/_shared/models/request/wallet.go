package request

type DefaultPaymentMethod struct {
	Number   string `json:"number"`
	ExpMonth string `json:"expMonth"`
	ExpYear  string `json:"expYear"`
	CVC      string `json:"cvc"`
}

type MakePaymentIntent struct {
	ProviderUserID        int    `json:"providerUserID"`
	AmountInSmallestDenom int64  `json:"amountInSmallestDenom"`
	Currency              string `json:"currency"`
	AutoCapture           bool   `json:"autoCapture"`
}

type ConfirmPaymentIntent struct {
	CustomerUserID int `json:"customerUserID"`
}

type MakeChatPaymentIntent struct {
	ProviderUserID int `json:"providerUserID"`
}

type GetTransactions struct {
	CursorID int64 `json:"cursorID"`
	Limit    int64 `json:"limit"`
}

type UpsertGoatChatsPrice struct {
	PriceInSmallestDenom int64  `json:"priceInSmallestDenom"`
	Currency             string `json:"currency"`
}

type UpsertBank struct {
	BankName          string         `json:"bankName" validate:"required|maxLength:200" message:"Please enter a valid bank name (max length 200)."`
	AccountNumber     string         `json:"accountNumber" validate:"required|maxLength:200|regex:^[0-9-]*$" message:"Please enter a valid account number (max length 200, only numbers)."`
	RoutingNumber     string         `json:"routingNumber" validate:"required|maxLength:200|regex:^[0-9-]*$" message:"Please enter a valid routing number (max length 200)."`
	AccountType       string         `json:"accountType" validate:"required|maxLength:50" message:"Please enter a valid type of bank account (e.g. checkings or savings) (max length 50)."`
	AccountHolderName string         `json:"accountHolderName" validate:"required|maxLength:200" message:"Please enter a valid account holder name"`
	AccountHolderType string         `json:"accountHolderType" validate:"required|maxLength:200" message:"Please enter a valid account holder type: (e.g. personal or business)"`
	Currency          string         `json:"currency" validate:"required|length:3" message:"Please enter a valid currency"`
	Country           string         `json:"country" validate:"required|maxLength:200" message:"Please enter a valid country"`
	BillingAddress    BillingAddress `json:"billingAddress"`
}

type BillingAddress struct {
	Street1    string `json:"street1" validate:"required|maxLength:200" message:"Please enter a valid street address (max length 200)."`
	Street2    string `json:"street2" validate:"maxLength:200" message:"Please enter a valid street address (max length 200)."`
	City       string `json:"city" validate:"required|maxLength:200" message:"Please enter a valid city (max length 200)."`
	State      string `json:"state" validate:"required|maxLength:200" message:"Please enter a valid state/region (max length 200)."`
	PostalCode string `json:"postalCode" validate:"required|maxLength:50" message:"Please enter a valid postal code (max length 50)."`
	Country    string `json:"country" validate:"required|maxLength:200" message:"Please enter a valid country (max length 200)."`
}

type GetBalance struct {
	Currency string `json:"currency"`
}

type GetGoatChatPrice struct {
	GoatUserID int `json:"goatUserID"`
}

type MarkProviderAsPaid struct {
	ProviderID  int   `json:"providerID"`
	AmountPaid  int64 `json:"amountPaid"`
	PayPeriodID int   `json:"payPeriodID"`
}

type ListUnpaidProviders struct{}

type UpsertBillingAddress struct {
	Street1    string `json:"street1" validate:"required|maxLength:200" message:"Please enter a valid street address (max length 200)."`
	Street2    string `json:"street2" validate:"maxLength:200" message:"Please enter a valid street address (max length 200)."`
	City       string `json:"city" validate:"required|maxLength:200" message:"Please enter a valid city (max length 200)."`
	State      string `json:"state" validate:"required|maxLength:200" message:"Please enter a valid state/region (max length 200)."`
	PostalCode string `json:"postalCode" validate:"required|maxLength:50" message:"Please enter a valid postal code (max length 50)."`
	Country    string `json:"country" validate:"required|maxLength:200" message:"Please enter a valid country (max length 200)."`
}

type GetPayoutPeriods struct {
	CursorID int64 `json:"cursorID"`
	Limit    int64 `json:"limit"`
}

type ListPayoutHistory struct {
	GoatUserID        int   `json:"goatUserID"`
	CursorPayPeriodID int64 `json:"cursorPayPeriodID"`
	Limit             int64 `json:"limit"`
}

type GetPayPeriod struct {
	Timestamp int64 `json:"timestamp"`
}
