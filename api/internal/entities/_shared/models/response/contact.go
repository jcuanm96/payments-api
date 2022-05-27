package response

type GetContacts struct {
	Paging   *Paging `json:"paging"`
	Contacts []User  `json:"contacts"`
}

type DeleteContact struct{}

type IsContact struct {
	IsContact bool `json:"isContact"`
}

type UploadContacts struct {
}

type ContactRecommendations struct {
	Recommendations []User `json:"recommendations"`
}
