package cloudtasks

const StripePaidGroupChatUnsubscribeQueueID string = "stripe-paid-group-unsubscribe"

type StripePaidGroupUnsubscribeTask struct {
	UserIDs    []int  `json:"userIDs"`
	ChannelID  string `json:"channelID"`
	GoatUserID int    `json:"goatUserID"`
}

const RemoveFromPaidGroupQueueID string = "remove-from-paid-group"

type RemoveFromPaidGroupTask struct {
	UserID    int    `json:"userID"`
	ChannelID string `json:"channelID"`
}

const AddUserContactsQueueID string = "add-user-contacts"

type AddUserContactsTask struct {
	UserID         int   `json:"userID"`
	ContactUserIDs []int `json:"contactUserIDs"`
}
