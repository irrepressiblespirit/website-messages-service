package entity

type MessageText struct {
	Message `bson:",inline"`
}

func (m *MessageText) ParseAddition(addition map[string]interface{}) error {
	return nil
}

type MessageImage struct {
	Message  `bson:",inline"`
	Addition struct {
		Size uint64 `json:"size" bson:"size"`
	} `json:"addition" bson:"addition"`
}

type MessageInRoomAnswer struct {
	Messages            []IMessage `json:"messages"`
	CompanionUreadCount int        `json:"companionunreadcount"`
}
