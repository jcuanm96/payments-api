package request

type GetBioData struct {
	Username string `json:"username"`
}

// This refers to a vama.me/username or vama.me/group-name-link
// whereas other Links in this file refer to Bio links for social medias.
type GetLink struct {
	LinkSuffix string `json:"linkSuffix"`
}

type NewMessageLink struct {
	MessageID string `json:"messageID"`
	ChannelID string `json:"channelID"`
}

type GetMessageByLink struct {
	LinkSuffix string `json:"linkSuffix"`
}

type GetThemes struct {
	CursorThemeID int `json:"cursorThemeID"`
	Limit         int `json:"limit"`
}

type UpsertBioLinks struct {
	TextContents []string `json:"textContents"`
	Links        []string `json:"links"`
	ThemeID      int      `json:"themeID"`
}

type UpsertTheme struct {
	Name                string `json:"name"`
	LogoColor           string `json:"logoColor"`
	IconColor           string `json:"iconColor"`
	TopGradientColor    string `json:"topGradientColor"`
	BottomGradientColor string `json:"bottomGradientColor"`
	UsernameColor       string `json:"usernameColor"`
	BioColor            string `json:"bioColor"`
	RowColor            string `json:"rowColor"`
	RowTextColor        string `json:"rowTextColor"`
}

type DeleteTheme struct {
	ThemeID int `json:"themeID"`
}
