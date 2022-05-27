package request

type Follow struct {
	UserID int `json:"userID"`
}

type Unfollow struct {
	UserID int `json:"userID"`
}

type IsFollowing struct {
	UserID int `json:"userID"`
}

type GetFollowedGoats struct {
	CursorGoatUserID int `json:"cursorGoatUserID"`
	Limit            int `json:"limit"`
}
