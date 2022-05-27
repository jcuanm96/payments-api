package appconfig

import (
	"log"

	"github.com/gopher-lib/config"
)

var Config AppConfig

func Init(filename string) {
	if err := config.LoadFile(&Config, filename); err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}
}

type Twilio struct {
	SID    string
	Token  string
	Verify string
}

type Stripe struct {
	Key         string
	EventSecret string
}

type Sendbird struct {
	MasterAPIKey  string
	ApplicationID string
}

type Gcloud struct {
	Project             string
	APIBaseURL          string
	RedirectBaseURL     string
	DynamicLinkBaseURL  string
	AppleBundleID       string
	ServiceAccountEmail string
}

type Postgre struct {
	ConnectionString string
}

type Storage struct {
	BucketName              string
	ProfileAvatarBucketName string
	FeedPostBucketName      string
	ThemesBucketName        string
	ChatMediaBucketName     string
}

type Redis struct {
	Endpoint string
	Port     string
	Password string
}

type AppConfig struct {
	Port           string
	Logger         bool
	ErrorReporting bool
	WebhookToken   string
	Gcloud         Gcloud     `mapstructure:"gcloud"`
	Twilio         Twilio     `mapstructure:"twilio"`
	Stripe         Stripe     `mapstructure:"stripe"`
	Sendbird       Sendbird   `mapstructure:"sendbird"`
	Postgre        Postgre    `mapstructure:"postgre"`
	Redis          Redis      `mapstructure:"redis"`
	Storage        Storage    `mapstructure:"storage"`
	Auth           Auth       `mapstructure:"auth"`
	Vote           Vote       `mapstructure:"vote"`
	Sharing        Sharing    `mapstructure:"sharing"`
	Telegram       Telegram   `mapstructure:"telegram"`
	VamaLogger     VamaLogger `mapstructure:"vamaLogger"`
}

type Auth struct {
	AccessTokenKey       string
	RefreshTokenKey      string
	GoatInviteCodeSecret string
	PayoutSecret         string
	DashboardSecret      string
}

type Vote struct {
	UpVote   string
	DownVote string
	Nil      string
}

type Sharing struct {
	ThemeSecret string
}

type Telegram struct {
	ApiToken   string
	ChannelIDs ChannelIDs
}

type ChannelIDs struct {
	VamaAlerts string
}

type VamaLogger struct {
	Type string
}
