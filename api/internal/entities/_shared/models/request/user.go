package request

import (
	"mime/multipart"

	"github.com/nyaruka/phonenumbers"
)

type UpdateUser struct {
	FirstName     *string `json:"firstName"`
	LastName      *string `json:"lastName"`
	Username      *string `json:"username"`
	Email         *string `json:"email"`
	ProfileAvatar *string `json:"profileAvatar"`
}

type ListUsers struct {
	Filters   string `json:"filters"`
	Query     string `json:"query"`
	UserTypes string `json:"userTypes"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}

type UploadProfileAvatar struct {
	ProfileAvatar *multipart.FileHeader `json:"profileAvatar"`
}

type UpsertBio struct {
	TextContent string `json:"textContent"`
}

type UpdatePhone struct {
	CountryCode string `json:"countryCode" query:"countryCode" message:"Please enter a valid country code."`
	Number      string `json:"phoneNumber" query:"phoneNumber" message:"Please enter a valid phone number."`
	Code        string `json:"code"`
}

func (p UpdatePhone) FormattedPhonenumber() (string, error) {
	num, formatErr := phonenumbers.Parse(p.Number, p.CountryCode)
	if formatErr != nil {
		return "", formatErr
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}

type GetUserContactByID struct {
	UserID int `json:"userID"`
}

type GetGoatProfile struct {
	GoatID int `json:"goatID"`
}

type GetGoatUsers struct {
	Limit       int  `json:"limit"`
	Page        int  `json:"page"`
	PageSize    int  `json:"pageSize"`
	ExcludeSelf bool `json:"excludeSelf"`
}

type BlockUser struct {
	UserID int `json:"userID"`
}

type UnblockUser struct {
	UserID int `json:"userID"`
}

type ReportUser struct {
	ReportedUserID int    `json:"reportedUserID"`
	Description    string `json:"description"`
}

type UpgradeUserToGoat struct {
	InviteCode string `json:"inviteCode"`
}
