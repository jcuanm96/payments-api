package request

type SearchGlobal struct {
	Query string `json:"query"`
}

type SearchMention struct {
	ChannelID string `json:"channelID"`
	Query     string `json:"query"`
}
