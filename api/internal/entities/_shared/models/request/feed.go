package request

import "mime/multipart"

type MakeFeedPost struct {
	TextContent        string                `json:"textContent"`
	Image              *multipart.FileHeader `json:"image"`
	GoatChatMessagesID int                   `json:"goatChatMessagesID"`
}

type MakeComment struct {
	Text   string `json:"text"`
	PostID int    `json:"postID"`
}

type DownvotePost struct {
	PostID int `json:"postID"`
}

type UpvotePost struct {
	PostID int `json:"postID"`
}

type GetUserFeedPosts struct {
	CursorPostID int `json:"cursorPostID"`
	Limit        int `json:"limit"`
}

type GetGoatFeedPosts struct {
	GoatUserID   int `json:"goatUserID"`
	CursorPostID int `json:"cursorPostID"`
	Limit        int `json:"limit"`
}

type GetFeedPostByID struct {
	PostID int `json:"postID"`
}

type GetFeedPostByLinkSuffix struct {
	LinkSuffix string `json:"linkSuffix"`
}

type DeleteFeedPost struct {
	PostID int `json:"postID"`
}

type GetFeedPostComments struct {
	PostID            int `json:"postID"`
	Limit             int `json:"limit"`
	LastUserCommentID int `json:"lastUserCommentID"`
}

type DeleteComment struct {
	CommentID int `json:"commentID"`
}

type ReactWebsocket struct {
	PostIDs []int `json:"postIDs"`
}
