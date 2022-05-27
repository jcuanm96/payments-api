package request

type UpdateFcmToken struct {
	Token string `json:"fcmToken"`
}

type SetGoatPostNotifications struct {
	GoatID int  `json:"goatID"`
	Enable bool `json:"enable"`
}

type UpdatePushSetting struct {
	ID         string `json:"id"`
	NewSetting string `json:"newSetting"`
}
