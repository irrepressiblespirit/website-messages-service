package entity

type MessageError struct {
	StatusCode  uint64            `json:"statusCode"`
	Message     string            `json:"message"`
	WrongFields map[string]string `json:"wrongFields"`
	MessageArgs []string          `json:"messageArgs"`
}

func (me MessageError) Error() string {
	return me.Message
}
