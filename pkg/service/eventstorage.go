package service

import "context"

type EventStorage interface {
	SendMessage(ctx context.Context, changel string, body []byte) error
}
