package response

type PushSetting struct {
	ID       string `json:"id"`
	Setting  string `json:"setting"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

type GetPushSettings struct {
	Settings []PushSetting `json:"settings"`
}
