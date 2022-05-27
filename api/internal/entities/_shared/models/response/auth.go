package response

type AuthConfirm struct{}

type AuthSuccess struct {
	Credentials Credentials `json:"credentials,omitempty"`
	User        User        `json:"user"`
}
type Credentials struct {
	AccessToken           string `json:"accessToken,omitempty"`
	AccessTokenExpiresAt  int64  `json:"accessTokenExpiresAt,omitempty"`
	RefreshToken          string `json:"refreshToken,omitempty"`
	RefreshTokenExpiresAt int64  `json:"refreshTokenExpiresAt,omitempty"`
	SendBirdAccessToken   string `json:"sendBirdAccessToken,omitempty"`
}

type Check struct {
	UserID  int    `json:"userID"`
	IsTaken bool   `json:"isTaken"`
	Message string `json:"message"`
}

type GenerateGoatInviteCode struct {
	GoatInviteCode string `json:"goatInviteCode"`
}

type CreateProviderAccount struct {
	URL       string `json:"url"`
	CreatedAt int64  `json:"createdAt"`
	ExpiresAt int64  `json:"expiresAt"`
	Object    string `json:"object"`
}

type VerifyProviderAccount struct {
	StripeAccountID     string `json:"stripeAccountID"`
	HasDetailsSubmitted bool   `json:"hasDetailsSubmitted"`
}
