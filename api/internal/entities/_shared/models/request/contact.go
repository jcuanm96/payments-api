package request

type CreateContact struct {
	ContactID int `json:"contactID"`
}

type DeleteContact struct {
	ContactID int `json:"contactID"`
}

type GetContacts struct {
	Page         int  `json:"page"`
	Size         int  `json:"size"`
	ExcludeGoats bool `json:"excludeGoats"`
}

type IsContact struct {
	ContactID int `json:"contactID"`
}

type Contact struct {
	Phone     string `json:"phoneNumber"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UploadContacts struct {
	Contacts []Contact `json:"contacts"`
}
