package response

import (
	"time"

	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

type (
	IsChatPublic struct {
		IsPublic bool `json:"isPublic"`
	}

	MessageRange struct {
		Start int    `json:"start"`
		End   int    `json:"end"`
		Type  string `json:"type"`
		ID    *int   `json:"id,omitempty"`
	}

	AdminMessageData struct {
		Ranges []MessageRange `json:"ranges"`
	}

	ChannelData struct {
		StartOfDraftMessageID *string `json:"startOfDraftMessageId,omitempty"`
		StartOfDraftMessageTS *int64  `json:"startOfDraftMessageTS,omitempty"`
		EndOfDraftMessageTS   *int64  `json:"endOfDraftMessageTS,omitempty"`
		ChannelState          *string `json:"channelState,omitempty"`
		IsConversationPublic  *bool   `json:"isConversationPublic,omitempty"`
		ExpiresAt             *int64  `json:"expiresAt,omitempty"`
	}

	GroupChannelData struct {
		LinkSuffix *string `json:"linkSuffix,omitempty"`
	}

	CreateGoatChatChannel struct {
		Channel *sendbird.GroupChannel `json:"channel"`
	}

	StartGoatChat struct {
		ChannelData ChannelData `json:"channelData"`
	}

	EndGoatChat struct {
		Post *FeedPost `json:"post"`
	}
	ConfirmGoatChat struct {
		Post        *FeedPost   `json:"post"`
		ChannelData ChannelData `json:"channelData"`
	}
	GetConversationEndTS struct {
		ConversationEndTS int64 `json:"conversationEndTS"`
	}

	ListChatMessages struct {
		Messages []sendbird.Message `json:"messages"`
	}

	GetGoatChatPostMessages struct {
		Messages []sendbird.Message `json:"messages"`
	}

	SendbirdUser struct {
		RequireAuthForProfileImage bool                 `json:"require_auth_for_profile_image"`
		IsActive                   bool                 `json:"is_active"`
		Role                       string               `json:"role"`
		UserID                     string               `json:"user_id"`
		Nickname                   string               `json:"nickname"`
		ProfileURL                 string               `json:"profile_url"`
		Metadata                   SendbirdUserMetadata `json:"metadata"`
	}

	SendbirdUserMetadata struct {
		Username string `json:"username"`
		UserType string `json:"userType"`
	}

	GetChannelWithUser struct {
		ChannelID string `json:"channelID"`
	}

	GetPaidGroup struct {
		GoatUser *User      `json:"goatUser"`
		Group    *PaidGroup `json:"group"`
	}

	CreatePaidGroupChannel struct {
		Channel sendbird.GroupChannel `json:"channel"`
	}

	LeavePaidGroup struct {
		Subscription *PaidGroupChatSubscription `json:"subscription"`
	}

	CancelPaidGroup struct {
		Subscription *PaidGroupChatSubscription `json:"subscription"`
	}

	GroupMetadata struct {
		Description *string `json:"description"`
		Benefit1    *string `json:"benefit1"`
		Benefit2    *string `json:"benefit2"`
		Benefit3    *string `json:"benefit3"`
	}

	PaidGroup struct {
		ID                   int                    `json:"id"`
		GoatID               int                    `json:"goatID"`
		PriceInSmallestDenom int64                  `json:"priceInSmallestDenom"`
		Currency             string                 `json:"currency"`
		ChannelID            string                 `json:"channelID"`
		LinkSuffix           string                 `json:"linkSuffix"`
		IsMember             bool                   `json:"isMember"`
		Metadata             GroupMetadata          `json:"metadata"`
		Channel              *sendbird.GroupChannel `json:"channel"`
		IsMemberLimitEnabled bool                   `json:"isMemberLimitEnabled"`
		MemberLimit          int                    `json:"memberLimit"`
	}

	ListGoatPaidGroups struct {
		GoatUser User        `json:"goatUser"`
		Groups   []PaidGroup `json:"groups"`
	}

	CheckLinkSuffixIsTaken struct {
		Taken bool `json:"taken"`
	}

	ListBannedUsers struct {
		Users []User `json:"users"`
	}

	DeepLinkInfo struct {
		PaidGroup     *PaidGroup `json:"paidGroup"`
		FreeGroup     *FreeGroup `json:"freeGroup"`
		UserChannelID *string    `json:"userChannelID"`
	}

	MediaURLs struct {
		UploadURL    string `json:"uploadURL"`
		ThumbnailURL string `json:"thumbnailURL"`
		ObjectID     string `json:"objectID"`

		DownloadURL          string `json:"downloadURL"`
		ThumbnailDownloadURL string `json:"thumbnailDownloadURL"`
	}

	MediaDownloadURL struct {
		ObjectID     string    `json:"objectID"`
		DownloadURL  string    `json:"downloadURL"`
		ThumbnailURL string    `json:"thumbnailURL"`
		ExpiresAt    time.Time `json:"expiresAt"`
	}

	BatchMediaDownloadURLs struct {
		URLs []MediaDownloadURL `json:"urls"`
	}

	VerifyMediaObject struct {
		Exists bool `json:"exists"`
	}

	CreateFreeGroupChannel struct {
		Channel sendbird.GroupChannel `json:"channel"`
	}

	FreeGroup struct {
		ID                   int                    `json:"id"`
		CreatedByUserID      int                    `json:"createdByUserID"`
		ChannelID            string                 `json:"channelID"`
		LinkSuffix           string                 `json:"linkSuffix"`
		IsMember             bool                   `json:"isMember"`
		Metadata             GroupMetadata          `json:"metadata"`
		Channel              *sendbird.GroupChannel `json:"channel"`
		IsMemberLimitEnabled bool                   `json:"isMemberLimitEnabled"`
		MemberLimit          int                    `json:"memberLimit"`
	}

	ListUserFreeGroups struct {
		User   User        `json:"user"`
		Groups []FreeGroup `json:"groups"`
	}

	GetFreeGroup struct {
		CreatedByUser *User      `json:"createdByUser"`
		Group         *FreeGroup `json:"group"`
	}
)
