package entity

import "fmt"

const (
	ErrInvalidToken           = "invalid token"
	ErrUnavailableType        = "unavailable message type"
	ErrRoomsCountLimit        = "rooms count limit"
	ErrRoomsUnsupportRoomType = "unsupport room type"
	ErrCoreMsgReqIDEquals     = "send message to himself"
	ErrMessagesCountLimit     = "count limit 100"
	ErrRefIDEquals            = "refid equals"
)

type CoreError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (ce CoreError) Error() string {
	return fmt.Sprintf("Error [%d] %s", ce.Code, ce.Message)
}

func CreateError(code int, text string) error {
	return CoreError{
		Code:    code,
		Message: text,
	}
}
