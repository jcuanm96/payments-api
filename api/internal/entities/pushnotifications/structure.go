package push

const (
	PENDING_BALANCE_REMINDERS_ID = "pending_balance"

	WALLET_CATEGORY = "Wallet"
)

// IDs must match column names in DB
var ValidSettingIDs = map[string]struct{}{
	PENDING_BALANCE_REMINDERS_ID: {},
}

type PushSettingConfig struct {
	ID       string
	Title    string
	Category string
}
type PushSettingConfigs struct {
	PendingBalance PushSettingConfig
}

var SettingConfigs = PushSettingConfigs{
	PendingBalance: PushSettingConfig{
		ID:       PENDING_BALANCE_REMINDERS_ID,
		Title:    "Pending Balance Reminders",
		Category: WALLET_CATEGORY,
	},
}

type UpdateSettings struct {
	PendingBalance string
}
