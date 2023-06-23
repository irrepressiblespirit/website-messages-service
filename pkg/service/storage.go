package service

import (
	"context"

	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
)

type Storage interface {
	// rooms
	GetRoom(ctx context.Context, roomid string) (*entity.Room, error)
	GetPrivateRoom(ctx context.Context, owner uint64, to uint64) (*entity.Room, error)
	GetUserRooms(ctx context.Context, user uint64, page int64, size int64) ([]*entity.Room, error)
	CreatePrivateRoom(ctx context.Context, room *entity.Room) (string, error)
	SetCurrentTimeInRoom(ctx context.Context, roomid string) error
	IncreaseUnreadCountInRoom(ctx context.Context, roomid string) error
	ZeroUnreadCountInRoom(ctx context.Context, roomid string, refid uint64) error
	GetUnreadCountMessagesByAllRooms(ctx context.Context, user uint64) (int64, error)

	// messages
	SaveMsg(ctx context.Context, msg entity.IMessage) (string, error)
	GetMsgs(ctx context.Context, roomid string, limit int64, lastMsgID string) ([]entity.IMessage, error)
	GetLastMsg(ctx context.Context, roomid string) (entity.IMessage, error)

	Ping(ctx context.Context) error
}

type StorageUsers interface {
	GetUser(ctx context.Context, refid uint64) (*entity.User, error)
	PutUser(ctx context.Context, user *entity.User) error
}
