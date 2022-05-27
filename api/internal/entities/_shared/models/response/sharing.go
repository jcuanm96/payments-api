package response

type GetBioData struct {
	UserID         int     `json:"userID"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	TextContent    *string `json:"textContent"`
	ProfileAvatar  *string `json:"profileAvatar"`
	NumFeedPosts   int     `json:"numFeedPosts"`
	NumGoatChats   int     `json:"numGoatChats"`
	NumSubscribers int     `json:"numSubscribers"`
}

type BioLinks struct {
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	ProfileAvatar *string  `json:"profileAvatar"`
	BioText       *string  `json:"bioText"`
	Username      string   `json:"username"`
	TextContents  []string `json:"textContents"`
	Links         []string `json:"links"`
	Theme         Theme    `json:"theme"`
}

type RedirectPaidGroupInfo struct {
	PriceInSmallestDenom int64         `json:"priceInSmallestDenom"`
	Currency             string        `json:"currency"`
	Metadata             GroupMetadata `json:"metadata"`
	MemberCount          int           `json:"memberCount"`
	Name                 string        `json:"name"`
	ProfileAvatar        string        `json:"profileAvatar"`
	Goat                 *User         `json:"goat"` // Partially filled
}

type RedirectFreeGroupInfo struct {
	Metadata      GroupMetadata `json:"metadata"`
	MemberCount   int           `json:"memberCount"`
	Name          string        `json:"name"`
	ProfileAvatar string        `json:"profileAvatar"`
	Goat          *User         `json:"goat"` // Partially filled, can be any user but named `goat` for backwards compatability with web code.
}

type GetVamaMeLink struct {
	BioData       *GetBioData            `json:"bioData"`
	PaidGroupInfo *RedirectPaidGroupInfo `json:"paidGroupInfo"`
	FreeGroupInfo *RedirectFreeGroupInfo `json:"freeGroupInfo"`
	DynamicLink   string                 `json:"dynamicLink"`
}

type UpsertBioLinks struct{}

type MessageLink struct {
	Link string `json:"link"`
}

type MessageInfo struct {
	ChannelID string `json:"channelID"`
	MessageID string `json:"messageID"`
}

type PublicMessage struct {
	DynamicLink string `json:"dynamicLink"`
	MessageText string `json:"messageText"`
	Sender      *User  `json:"sender"` // Partially filled
}

type GetThemes struct {
	Themes []Theme `json:"themes"`
}

type Theme struct {
	ID                  int    `json:"id"`
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

type UpsertTheme struct{}
