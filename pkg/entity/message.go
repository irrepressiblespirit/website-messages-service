package entity

import (
	"strings"
	"time"
)

type MessageType string

const (
	TypeText  = "text"
	TypeFile  = "file"
	TypeImage = "image"
	TypeURL   = "url"
)

func GetMessageStructByType(stringType MessageType) (IMessage, error) {
	switch stringType {
	case TypeText:
		return &MessageText{}, nil
	case TypeFile:
		return &MessageFile{}, nil
	default:
		return nil, CreateError(500, ErrUnavailableType)
	}
}

const CurrentMessageVersion uint64 = 1

type Message struct {
	Version uint64      `json:"version" bson:"version"`
	ID      string      `json:"id,omitempty" bson:"_id,omitempty"`
	RoomID  string      `json:"roomid" bson:"roomid"`
	RefID   uint64      `json:"refid" bson:"refid"`
	Sended  time.Time   `json:"sended" bson:"sended"`
	Type    MessageType `json:"type" bson:"type"`
	Body    string      `json:"body" bson:"body"`
}

type IMessage interface {
	SetVersion(uint64) IMessage
	GetVersion() uint64
	GetID() string
	SetID(string) IMessage
	GetRoomID() string
	SetRoomID(string) IMessage
	GetRefID() uint64
	SetRefID(uint64) IMessage
	SetSended(time.Time) IMessage
	GetSended() time.Time
	GetType() MessageType
	SetType(MessageType) IMessage
	GetBody() string
	SetBody(string) IMessage
	ParseAddition(map[string]interface{}) error
	GetAddition() map[string]interface{}
	IsCorrect() error
}

func (m *Message) GetVersion() uint64                         { return m.Version }
func (m *Message) GetID() string                              { return m.ID }
func (m *Message) GetRoomID() string                          { return m.RoomID }
func (m *Message) GetRefID() uint64                           { return m.RefID }
func (m *Message) GetSended() time.Time                       { return m.Sended }
func (m *Message) GetType() MessageType                       { return m.Type }
func (m *Message) GetBody() string                            { return m.Body }
func (m *Message) ParseAddition(map[string]interface{}) error { return nil }
func (m *Message) GetAddition() map[string]interface{}        { return make(map[string]interface{}) }

func (m *Message) IsCorrect() error {
	if m.Body == "" || len(strings.TrimSpace(m.Body)) == 0 {
		return MessageError{
			StatusCode: 3,
			Message:    "errors.messages.body.is.empty",
		}
	}
	return nil
}

func (m *Message) SetVersion(version uint64) IMessage {
	m.Version = version
	return m
}

func (m *Message) SetID(id string) IMessage {
	m.ID = id
	return m
}

func (m *Message) SetRoomID(roomid string) IMessage {
	m.RoomID = roomid
	return m
}

func (m *Message) SetRefID(refid uint64) IMessage {
	m.RefID = refid
	return m
}

func (m *Message) SetSended(t time.Time) IMessage {
	m.Sended = t
	return m
}

func (m *Message) SetType(mt MessageType) IMessage {
	m.Type = mt
	return m
}

func (m *Message) SetBody(body string) IMessage {
	m.Body = body
	return m
}
