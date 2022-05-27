package request

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

type SignupSMS struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone
	Email string `json:"email"`
	Code  string `json:"code"` // Twilio code
}

type SignupSMSGoat struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone
	Username   string `json:"username"`
	Email      string `json:"email"`
	InviteCode string `json:"inviteCode"`
	Code       string `json:"code"` // Twilio code
}

type VerifySMS struct {
	Phone
}

type SignInSMS struct {
	Phone
	Code string `json:"code"`
}

type Username struct {
	Username string `json:"username" query:"username"`
}

type SendBirdUserMetadata struct {
	Username
	Type string `json:"userType"`
}

type Email struct {
	Email string `json:"email" query:"email"`
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken"`
}

type SignOut struct {
	RefreshToken string `json:"refreshToken"`
}

type GenerateGoatInviteCode struct{}

type GoatInviteCode struct {
	Code string `json:"code"`
}

type CreateProviderAccount struct {
	StripeAccountID      *string `json:"stripeAccountID"`
	GoatInviteCodeSecret string  `json:"goatInviteCodeSecret"`
	Email                string  `json:"email"`
}

type VerifyProviderAccount struct {
	StripeAccountID *string `json:"stripeAccountID"`
	CountryCode     string  `json:"countryCode" query:"countryCode"`
	Number          string  `json:"phoneNumber" query:"phoneNumber"`
}

type Phone struct {
	CountryCode string `json:"countryCode" query:"countryCode"`
	Number      string `json:"phoneNumber" query:"phoneNumber"`
}

func (p Phone) Validate() bool {
	num, err := phonenumbers.Parse(p.Number, p.CountryCode)
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(num)
}

// formattedPhonenumber returns formatted phone number if .Number is valid phone number,
// otherwise returns empty string.
func (p Phone) FormattedPhonenumber() (string, error) {
	num, err := phonenumbers.Parse(p.Number, p.CountryCode)
	if err != nil {
		return "", err
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}

func (p Phone) HidePartPhoneSymbols() (string, error) {
	s, err := p.FormattedPhonenumber()
	if err != nil {
		return "", err
	}
	res := fmt.Sprintf("%s *** *** ** %s", string(s[0:2]), string(s[len(s)-2:]))
	return res, nil
}
