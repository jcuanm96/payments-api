package sendbird

// https://sendbird.com/docs/chat/v3/platform-api/guides/webhooks
const GroupChannelCreateCategory = "group_channel:create"
const GroupChannelChangedCategory = "group_channel:changed"
const GroupChannelInviteCategory = "group_channel:invite"
const GroupChannelJoinCategory = "group_channel:join"
const GroupChannelLeaveCategory = "group_channel:leave"
const GroupChannelMessageSendCategory = "group_channel:message_send"
const GroupChannelMessageUpdateCategory = "group_channel:message_update"
const GroupChannelMessageDeleteCategory = "group_channel:message_delete"

// Generic event to unmarshal the category to then switch
// on the actual event struct to use.
type Event struct {
	Category string `json:"category"`
}

type WebhookUser struct {
	UserID     string `json:"user_id"`
	NickName   string `json:"nickname"`
	ProfileURL string `json:"profile_url"`
	Inviter    User   `json:"inviter"` // this still only has WebhookUser fields
}

type GroupChannelCreateEvent struct {
	Category  string       `json:"category"`
	CreatedAt int64        `json:"created_at"`
	Members   []User       `json:"members"`
	Inviter   WebhookUser  `json:"inviter"`
	Channel   GroupChannel `json:"channel"`
	AppID     string       `json:"app_id"`
}

type GroupChannelChange struct {
	Key string `json:"key"`
	Old string `json:"old"`
	New string `json:"new"`
}

type GroupChannelChangedEvent struct {
	Category  string               `json:"category"`
	CreatedAt int64                `json:"created_at"`
	Changes   []GroupChannelChange `json:"changes"`
	Members   []User               `json:"members"`
	Channel   GroupChannel         `json:"channel"`
	AppID     string               `json:"app_id"`
}

type GroupChannelInviteEvent struct {
	Category  string        `json:"category"`
	InvitedAt int64         `json:"invited_at"`
	Members   []User        `json:"members"`
	Inviter   WebhookUser   `json:"inviter"`
	Channel   GroupChannel  `json:"channel"`
	Invitees  []WebhookUser `json:"invitees"`
	AppID     string        `json:"app_id"`
}

type GroupChannelJoinEvent struct {
	Category string        `json:"category"`
	JoinedAt int64         `json:"joined_at"`
	Members  []User        `json:"members"`
	Channel  GroupChannel  `json:"channel"`
	Users    []WebhookUser `json:"users"` // Users who have joined the channel
	AppID    string        `json:"app_id"`
}
