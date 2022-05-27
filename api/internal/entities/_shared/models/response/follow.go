package response

type IsFollowing struct {
	IsFollowing bool `json:"isFollowing"`
}

type GetFollowedGoats struct {
	Goats []User `json:"goats"`
}
