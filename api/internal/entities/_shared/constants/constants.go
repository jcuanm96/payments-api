package constants

const (
	DEFAULT_TIER_NAME = "TIERLESS"
	WITHDRAW          = "WITHDRAW"
	PERSONAL          = "PERSONAL"
	BUSINESS          = "BUSINESS"

	TXN_HISTORY_INCOMING           = "INCOMING"
	TXN_HISTORY_OUTGOING           = "OUTGOING"
	TXN_HISTORY_PAYOUT             = "PAYOUT"
	TXN_HISTORY_PENDING            = "PENDING"
	TXN_HISTORY_SOURCE_TYPE_CHARGE = "charge"

	TOTAL_FEES_RATIO = float64(0.05)

	VAMA_WEB_BASE_URL = "https://vama.com"

	// strings to be passed to fmt.Sprintf or similar
	GCS_URL_F            = "https://storage.googleapis.com/%s/%s"
	FEED_POST_BASE_URL_F = "%s/p/%s"
	MESSAGE_BASE_URL_F   = "%s/m/%s"

	MAX_POST_PREVIEW_CONTENT_LENGTH = 300

	PUSH_NOTIFICATION_ON    = "ON"
	PUSH_NOTIFICATION_OFF   = "OFF"
	PUSH_NOTIFICATION_UNSET = "UNSET"

	// Sendbird channel data constants
	ChannelStatePending   = "PENDING"
	ChannelStateActive    = "ACTIVE"
	ChannelTypePaidDirect = "PAID_DIRECT"
	ChannelTypeDirect     = "DIRECT"
	ChannelTypePaidGroup  = "PAID_GROUP"
	ChannelTypeFreeGroup  = "FREE_GROUP"

	ChannelEventCustomType         = "CHANNEL_EVENT"
	ChannelEventRemoveCustomType   = "CHANNEL_EVENT_REMOVE"
	PendingBalanceNotificationType = "PENDING_BALANCE"

	AdminMessageRangeNickname = "NICKNAME"
	AdminMessageRangeText     = "TEXT"

	MIN_REQUIRED_IOS_APP_VERSION = "1.0.0"
	API_V2                       = "/api/v2"
	PUBLIC_V1                    = "/public"
	WEBSOCKET_V1                 = "/ws/v1"
	CLOUD_TASK_V1                = "/cloudtask/v1"
	CLOUD_SCHEDULER_V1           = "/cloudscheduler/v1"

	GOAT_CHAT_EXPIRES_AT_DEFAULT_VALUE = -1

	ACCESS_TOKEN_TYPE  = "AccessToken"
	REFRESH_TOKEN_TYPE = "RefreshToken"

	VAMA_USER_TYPE = "VAMA"
	VAMA_USER_ID   = -1

	VAMA_THEME_NAME = "Vama"
	MAX_LINKS       = 7

	DEFAULT_CURRENCY = "usd"

	MIN_GOAT_CHAT_PRICE_USD_IN_SMALLEST_DENOM = 300
	MAX_GOAT_CHAT_PRICE_USD_IN_SMALLEST_DENOM = 10000

	MIN_PAID_GROUP_CHAT_PRICE_USD_IN_SMALLEST_DENOM = 300
	MAX_PAID_GROUP_CHAT_PRICE_USD_IN_SMALLEST_DENOM = 100000

	MIN_USER_SUBSCRIPTION_PRICE_USD_IN_SMALLEST_DENOM = 300
	MAX_USER_SUBSCRIPTION_PRICE_USD_IN_SMALLEST_DENOM = 10000

	GOAT_SIGNUP_CREDIT_AMOUNT_USD_IN_SMALLEST_DENOM = 10000

	MAX_GROUP_NAME_LENGTH  = 140
	MAX_LINK_SUFFIX_LENGTH = 110

	MAX_UPLOAD_CONTACT_NAME_LENGTH = 50
	MAX_NAME_LENGTH_FOR_USER       = 50

	MAX_USER_BIO_LENGTH                = 1000
	MAX_REPORT_USER_DESCRIPTION_LENGTH = 1000

	MAX_SENDBIRD_CHANNEL_ID_LEN = 100

	MAX_GROUP_DESCRIPTION_LENGTH = 500
	MAX_GROUP_BENEFIT_LENGTH     = 130

	DEFAULT_PAID_GROUP_MEMBER_LIMIT            = 100
	MIN_GROUP_MEMBER_LIMIT                     = 2
	MAX_GROUP_MEMBER_LIMIT                     = 2000
	DEFAULT_PAID_GROUP_IS_MEMBER_LIMIT_ENABLED = false

	GCP_LOGGER_TYPE    = "GCP"
	APPLE_APP_STORE_ID = "1600580466"

	CURR_PAYMENTS_VERSION = "stripe_v2"
)

var CreatorDiscoverPriority = map[string][]int{
	"vama-dev":     []int{2},
	"vama-staging": []int{18},
	"vama-prod":    []int{72, 292},
	"vama-test":    []int{},
}
