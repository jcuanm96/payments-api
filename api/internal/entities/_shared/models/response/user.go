package response

import (
	"time"
)

type User struct {
	ID              int        `json:"id"`
	UUID            string     `json:"uuid"`
	StripeID        *string    `json:"-"`
	FirstName       string     `json:"firstName"`
	LastName        string     `json:"lastName"`
	Phonenumber     string     `json:"phoneNumber"`
	CountryCode     string     `json:"countryCode"`
	Email           *string    `json:"email"`
	Username        string     `json:"username,omitempty"`
	Type            string     `json:"userType"`
	ProfileAvatar   *string    `json:"profileAvatar"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty"`
	StripeAccountID *string    `json:"-"` // deprecated
}

type ListUsers struct {
	Paging Paging `json:"paging"`
	Users  []User `json:"users"`
}

type UploadProfileAvatar struct {
	Filename string `json:"filename"`
}

type UpsertBio struct{}

type UpdatePhone struct{}

type SoftDeleteCurrentUser struct{}

type GetGoatUsers struct {
	GoatUsers []GetGoatUsersListItem `json:"goatUsers"`
}

type GetGoatUsersListItem struct {
	User         User `json:"user"`
	NumGoatChats int  `json:"numGoatChats"`
}

type GetUserContactByID struct {
	User      *User `json:"user"`
	IsContact bool  `json:"isContact"`
	IsBlocked bool  `json:"isBlocked"`
}

type GetGoatProfile struct {
	User                     User `json:"user"`
	IsFollowing              bool `json:"isFollowing"`
	PostNotificationsEnabled bool `json:"postNotificationsEnabled"`
}

type Invite struct {
	Code string `json:"code"`
	User *User  `json:"user"`
}

type MyProfile struct {
	User User   `json:"user"`
	Bio  string `json:"bio"`
}

type GetInviteCodeStatuses struct {
	Invites []Invite `json:"invites"`
}
