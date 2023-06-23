package entity

import "time"

const (
	RoomTypePrivate = "private"
	RoomTypeGroup   = "group"
)

type Room struct {
	ID          string      `json:"id,omitempty" bson:"_id,omitempty"`
	Type        string      `json:"type" bson:"type"`
	Name        string      `json:"name" bson:"name"` // используется только в групповом (название комнаты)
	LastMessage time.Time   `json:"last_msg" bson:"last_msg"`
	Users       []RoomUsers `json:"users" bson:"users"`
}

type RoomUsers struct {
	RefID        uint64 `json:"refid" bson:"refid"`
	Admin        bool   `json:"admin" bson:"admin"`     // для групповых чатов
	NotReadCount int    `json:"notread" bson:"notread"` // сколько у этого пользователя в этой комнате не прочитанно сообщений
}

type RoomNotFoundError struct{}

func (r RoomNotFoundError) Error() string {
	return "Room not found"
}

type MyRoom struct {
	ID                   string         `json:"id"`
	Type                 string         `json:"type"`
	Name                 string         `json:"name"`
	LastMessage          time.Time      `json:"lastmessage"`
	LogoURL              string         `json:"logourl"`
	NotReadCount         int            `json:"notread"`
	CompanionUnreadCount int            `json:"companionunreadcount"`
	CompanionRefID       uint64         `json:"companionrefid"`
	CompanionOwnerRefID  uint64         `json:"companionownerrefid"`
	LastMsg              *MyRoomLastMsg `json:"last"`
}

type MyRoomLastMsg struct {
	ID    string      `json:"id"`
	RefID uint64      `json:"refid"`
	Type  MessageType `json:"type"`
	Body  string      `json:"body"`
}

type MyRoomWithNotReadMessagesTotalCount struct {
	Room                 MyRoom `json:"room"`
	NotReadMessagesCount int64  `json:"notreadmessagescount"`
}

type RoomAccessDeniedError struct{}

func (r RoomAccessDeniedError) Error() string {
	return "Access Denied"
}

func (room *Room) UserExist(refid uint64) bool {
	for _, v := range room.Users {
		if v.RefID == refid {
			return true
		}
	}
	return false
}

func (room *Room) GetUser(refid uint64) *RoomUsers {
	for i := 0; i < len(room.Users); i++ {
		if room.Users[i].RefID == refid {
			return &room.Users[i]
		}
	}
	return nil
}

func (room *Room) GetOtherUser(refid uint64) *RoomUsers {
	if room.Type == RoomTypePrivate {
		for i := 0; i < len(room.Users); i++ {
			if room.Users[i].RefID != refid {
				return &room.Users[i]
			}
		}
	}
	return nil
}
