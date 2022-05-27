package request

import "mime/multipart"

type (
	IsChatPublic struct {
		SendBirdChannelID string `json:"sendBirdChannelID"`
	}

	CreateGoatChatChannel struct {
		OtherUserID int `json:"userID"`
	}

	StartGoatChat struct {
		SendBirdChannelID string `json:"sendBirdChannelID"`
		IsPublic          bool   `json:"isPublic"`
	}

	EndGoatChat struct {
		SendBirdChannelID string `json:"sendBirdChannelID"`
		StartTS           int64  `json:"startTS"`
		IsPublic          bool   `json:"isPublic"`
	}

	ConfirmGoatChat struct {
		SendBirdChannelID string `json:"sendBirdChannelID"`
		CustomerUserID    int    `json:"customerUserID"`
		IsPublic          bool   `json:"isPublic"`
	}

	GetConversationEndTS struct {
		SendBirdChannelID string `json:"sendBirdChannelID"`
	}

	ListChatMessages struct {
		SendBirdChannelID string `json:"sendBirdChannelID"`
		MessageTsFrom     int64  `json:"messageTsFrom"`
		MessageTsTo       *int64 `json:"messageTsTo"`
		Limit             int    `json:"limit"`
	}

	GetGoatChatPostMessages struct {
		PostID           int   `json:"postID"`
		Limit            int   `json:"limit"`
		CursorMesssageTS int64 `json:"cursorMessageTS"`
	}

	GetChannelWithUser struct {
		UserID int `json:"userID"`
	}

	JoinPaidGroup struct {
		ChannelID string `json:"channelID"`
	}

	GetPaidGroup struct {
		ChannelID string `json:"channelID"`
	}

	CreatePaidGroupChannel struct {
		PriceInSmallestDenom int                   `json:"priceInSmallestDenom"`
		Currency             string                `json:"currency"`
		Name                 string                `json:"name"`
		CoverFile            *multipart.FileHeader `json:"coverFile"`
		LinkSuffix           string                `json:"linkSuffix"`

		Description          *string `json:"description"`
		Benefit1             *string `json:"benefit1"`
		Benefit2             *string `json:"benefit2"`
		Benefit3             *string `json:"benefit3"`
		Members              []int   `json:"members"`
		MemberLimit          *int    `json:"memberLimit"`
		IsMemberLimitEnabled *bool   `json:"isMemberLimitEnabled"`
	}

	UpdatePaidGroupPrice struct {
		ChannelID            string `json:"channelID"`
		PriceInSmallestDenom int    `json:"priceInSmallestDenom"`
		Currency             string `json:"currency"`
	}

	CheckLinkSuffixIsTaken struct {
		LinkSuffix string  `json:"linkSuffix"`
		ChannelID  *string `json:"channelID"`
	}

	UpdateGroupLink struct {
		ChannelID  string `json:"channelID"`
		LinkSuffix string `json:"linkSuffix"`
	}

	UpdateGroupMetadata struct {
		ChannelID   string  `json:"channelID"`
		Description *string `json:"description"`
		Benefit1    *string `json:"benefit1"`
		Benefit2    *string `json:"benefit2"`
		Benefit3    *string `json:"benefit3"`
	}

	LeaveGroup struct {
		ChannelID string `json:"channelID"`
	}

	CancelPaidGroup struct {
		ChannelID string `json:"channelID"`
	}

	ListGoatPaidGroups struct {
		GoatID   int   `json:"goatID"`
		CursorID int   `json:"cursorID"`
		Limit    int64 `json:"limit"`
	}

	DeletePaidGroup struct {
		ChannelID string `json:"channelID"`
	}

	BanUserFromPaidGroup struct {
		BannedUserID int    `json:"userID" query:"userID"`
		ChannelID    string `json:"channelID"`
	}

	UnbanUserFromPaidGroup struct {
		BannedUserID int    `json:"userID"`
		ChannelID    string `json:"channelID"`
	}

	RemoveUserFromGroup struct {
		UserID    int    `json:"userID"`
		ChannelID string `json:"channelID"`
	}

	ListBannedUsers struct {
		ChannelID string `json:"channelID"`
		CursorID  int    `json:"cursorID"`
		Limit     int64  `json:"limit"`
	}

	UpdatePaidGroupMemberLimits struct {
		ChannelID            string `json:"channelID"`
		MemberLimit          *int   `json:"memberLimit"`
		IsMemberLimitEnabled *bool  `json:"isMemberLimitEnabled"`
	}

	GetDeepLinkInfo struct {
		LinkSuffix string `json:"linkSuffix"`
	}

	GetBatchMediaDownloadURLs struct {
		ObjectIDs []string `json:"objectIDs"`
	}

	VerifyMediaObject struct {
		ObjectID string `json:"objectID"`
	}

	CreateFreeGroupChannel struct {
		Name       string                `json:"name"`
		CoverFile  *multipart.FileHeader `json:"coverFile"`
		LinkSuffix string                `json:"linkSuffix"`

		Description          *string `json:"description"`
		Benefit1             *string `json:"benefit1"`
		Benefit2             *string `json:"benefit2"`
		Benefit3             *string `json:"benefit3"`
		Members              []int   `json:"members"`
		MemberLimit          *int    `json:"memberLimit"`
		IsMemberLimitEnabled *bool   `json:"isMemberLimitEnabled"`
	}

	JoinFreeGroup struct {
		ChannelID string `json:"channelID"`
	}

	GetFreeGroup struct {
		ChannelID string `json:"channelID"`
	}

	ListUserFreeGroups struct {
		UserID   int   `json:"userID"`
		CursorID int   `json:"cursorID"`
		Limit    int64 `json:"limit"`
	}

	DeleteFreeGroup struct {
		ChannelID string `json:"channelID"`
	}

	AddFreeGroupChatCoCreators struct {
		ChannelID  string `json:"channelID"`
		CoCreators []int  `json:"coCreators"`
	}

	RemoveFreeGroupChatCoCreator struct {
		UserID    string `json:"userID"`
		ChannelID string `json:"channelID"`
	}
)
