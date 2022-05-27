package response

import (
	"time"

	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

type MakeFeedPost struct {
	Creator User `json:"creator"`
}

type Reaction struct {
	NewState     string `json:"newState"`
	NumUpvotes   int64  `json:"numUpvotes"`
	NumDownvotes int64  `json:"numDownvotes"`
}

type MakeComment struct{}

type GetUserFeedPosts struct {
	FeedPosts []FeedPost `json:"feedPosts"`
}

type GetGoatFeedPosts struct {
	FeedPosts []FeedPost `json:"feedPosts"`
}

type PostImage struct {
	URL string `json:"url"`
	// width/height in pixels
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Conversation struct {
	ChannelID       string             `json:"channelID"`
	StartTS         uint64             `json:"startTS"`
	EndTS           uint64             `json:"endTS"`
	PreviewMessages []sendbird.Message `json:"previewMessages"`
}

type FeedPost struct {
	ID               int          `json:"id"`
	UserID           int          `json:"userID"`
	PostCreatedAt    time.Time    `json:"postCreatedAt"`
	PostNumUpvotes   int64        `json:"postNumUpvotes"`
	PostNumDownvotes int64        `json:"postNumDownvotes"`
	PostTextContent  string       `json:"postTextContent"`
	Conversation     Conversation `json:"conversation"`
	NumComments      int          `json:"numComments"`
	Image            PostImage    `json:"image"`
	Reaction         string       `json:"reaction"`
	Creator          User         `json:"creator"`
	Customer         *User        `json:"customer"`
	Link             string       `json:"link"`
}

type PublicFeedPost struct {
	PostCreatedAt    time.Time `json:"postCreatedAt"`
	PostNumUpvotes   int64     `json:"postNumUpvotes"`
	PostNumDownvotes int64     `json:"postNumDownvotes"`
	PostTextContent  string    `json:"postTextContent"`
	IsConversation   bool      `json:"isConversation"`
	NumComments      int       `json:"numComments"`
	Image            PostImage `json:"image"`
	Creator          User      `json:"creator"`  // Only partially filled
	Customer         *User     `json:"customer"` // Only partially filled
	Link             string    `json:"link"`
	DynamicLink      string    `json:"dynamicLink"`
}

type DeleteFeedPost struct{}

type GetFeedPostComments struct {
	Comments []Comment `json:"comments"`
}

type CommentMetadata struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userID"`
	PostID      int       `json:"postID"`
	TextContent string    `json:"textContent"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt"`
}

type Comment struct {
	CommentMetadata `json:"comment"`
	User            User `json:"user"`
}
