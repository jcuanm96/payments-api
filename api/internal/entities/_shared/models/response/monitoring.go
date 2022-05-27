package response

type GetMinRequiredVersionByDevice struct {
	Version string `json:"version"`
}

type GetDashboard struct {
	Total DashboardStats `json:"total"`
	Day   DashboardStats `json:"day"`
	Week  DashboardStats `json:"week"`
	Month DashboardStats `json:"month"`
	Year  DashboardStats `json:"year"`
}

type DashboardStats struct {
	NewUsers  int `json:"users"`
	GoatChats int `json:"goatChats"`
	FeedPosts int `json:"feedPosts"`
}
