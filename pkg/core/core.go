package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
	"github.com/irrepressiblespirit/website-messages-service/pkg/service"
)

const (
	SLAMaxGetMessages = 100
	SLAMaxRoomsLimit  = 100
)

type Core struct {
	Storage      service.Storage
	StorageUsers service.StorageUsers
	TokenService service.TokenService
	EventStorage service.EventStorage
	AuthToken    service.IAuthToken
	Config       *entity.Config
}

func (c *Core) check() error {
	if c.Storage == nil {
		return entity.CreateError(0, "Storage not configured")
	}
	if c.StorageUsers == nil {
		return entity.CreateError(0, "StorageUsers not configured")
	}
	if c.TokenService == nil {
		return entity.CreateError(0, "TokenService not configured")
	}
	if c.EventStorage == nil {
		return entity.CreateError(0, "EventStorage not configured")
	}
	if c.AuthToken == nil {
		return entity.CreateError(0, "AuthToken not configured")
	}
	return nil
}

// /////////////// GetUser logic ////
/*
get user from local storage
return user if his present, if not send request to external user service
*/
func (c *Core) GetUser(ctx context.Context, refid uint64) (*entity.User, error) {
	if err := c.check(); err != nil {
		return nil, err
	}
	return c.StorageUsers.GetUser(ctx, refid)
}

func (c *Core) GetNewToken(refid uint64) (entity.Token, error) {
	return c.TokenService.Generate(refid)
}

// return list of rooms with not read messages count
func (c *Core) GetMyRooms(ctx context.Context, refid uint64, page int64, size int64) ([]entity.MyRoom, int64, error) {
	if size > SLAMaxRoomsLimit || size <= 0 {
		return nil, 0, entity.CreateError(400, entity.ErrRoomsCountLimit)
	}
	result := []entity.MyRoom{}
	rooms, err := c.Storage.GetUserRooms(ctx, refid, page, size)
	if err != nil {
		return result, 0, err
	}
	for _, item := range rooms {
		switch item.Type {
		case entity.RoomTypePrivate:
			myroom, err := c.ConvertPrivateRoomToMyRoom(ctx, item, refid)
			if err != nil {
				return nil, 0, err
			}
			result = append(result, *myroom)
		default:
			return result, 0, entity.CreateError(500, entity.ErrRoomsUnsupportRoomType)
		}
	}
	count, _ := c.GetUnreadCount(ctx, refid)
	return result, count, nil
}

func (c *Core) GetUnreadCount(ctx context.Context, refid uint64) (int64, error) {
	count, err := c.Storage.GetUnreadCountMessagesByAllRooms(ctx, refid)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *Core) GetPrivateRoom(ctx context.Context, user1 uint64, user2 uint64) (*entity.Room, error) {
	if user1 == user2 {
		return nil, entity.CreateError(409, entity.ErrCoreMsgReqIDEquals)
	}
	room, err := c.Storage.GetPrivateRoom(ctx, user1, user2)
	if err == nil {
		return room, nil
	}
	if errors.Is(err, entity.RoomNotFoundError{}) {
		_, err = c.GetUser(ctx, user1)
		if err != nil {
			return nil, err
		}
		_, err = c.GetUser(ctx, user2)
		if err != nil {
			return nil, err
		}
		room = &entity.Room{
			Type:        entity.RoomTypePrivate,
			Name:        "none",
			LastMessage: time.Now(),
			Users: []entity.RoomUsers{
				{RefID: user1, Admin: false, NotReadCount: 0},
				{RefID: user2, Admin: false, NotReadCount: 0},
			},
		}
		newID, err := c.Storage.CreatePrivateRoom(ctx, room)
		if err != nil {
			return nil, err
		}
		room.ID = newID
		return room, nil
	}
	return nil, err
}

func (c *Core) SendMsg(ctx context.Context, msg entity.IMessage) (string, error) {
	room, err := c.Storage.GetRoom(ctx, msg.GetRoomID())
	if err != nil {
		return "", fmt.Errorf("error found room: %w", err)
	}
	if !room.UserExist(msg.GetRefID()) {
		return "", entity.RoomAccessDeniedError{}
	}
	msg.SetSended(time.Now())
	if err := msg.IsCorrect(); err != nil {
		log.Printf("error check correct: %+v", err)
		return "", err
	}
	newID, err := c.Storage.SaveMsg(ctx, msg)
	if err != nil {
		log.Printf("error save msg: %+v", err)
		return "", err
	}
	err = c.Storage.SetCurrentTimeInRoom(ctx, room.ID)
	if err != nil {
		return "", err
	}
	err = c.Storage.IncreaseUnreadCountInRoom(ctx, room.ID)
	if err != nil {
		return "", err
	}
	err = c.Storage.ZeroUnreadCountInRoom(ctx, room.ID, msg.GetRefID())
	if err != nil {
		return "", err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := c.sendMyRoom(ctx, room.ID)
		if err != nil {
			log.Printf("Error: core.SendFastMsg sendMyRoom %+v", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := c.sendMessageInRoom(ctx, room, msg)
		if err != nil {
			log.Printf("Error core.SendFastMsg sendMessageInRoom %+v", err)
		}
	}()
	wg.Wait()
	return newID, nil
}

func (c *Core) SendFastMsg(ctx context.Context, from uint64, to uint64, body string) (string, string, error) {
	room, err := c.GetPrivateRoom(ctx, from, to)
	if err != nil {
		return "", "", err
	}
	msg := entity.MessageText{
		Message: entity.Message{
			Version: entity.CurrentMessageVersion,
			RoomID:  room.ID,
			RefID:   from,
			Type:    entity.TypeText,
			Body:    body,
		},
	}
	messageID, err := c.SendMsg(ctx, &msg)
	return room.ID, messageID, err
}

func (c *Core) GetMessages(ctx context.Context, roomid string, refid uint64, count int, lastMsgID string) (*entity.MessageInRoomAnswer, error) {
	if count > SLAMaxGetMessages {
		return nil, entity.CreateError(400, entity.ErrMessagesCountLimit)
	}
	room, err := c.Storage.GetRoom(ctx, roomid)
	if err != nil {
		return nil, err
	}
	if !room.UserExist(refid) {
		return nil, entity.RoomAccessDeniedError{}
	}
	msgs, err := c.Storage.GetMsgs(ctx, roomid, int64(count), lastMsgID)
	if err != nil {
		return nil, err
	}
	CompanionUreadCount := 0
	if room.Type == entity.RoomTypePrivate {
		CompanionUreadCount = room.GetOtherUser(refid).NotReadCount
	}
	return &entity.MessageInRoomAnswer{Messages: msgs, CompanionUreadCount: CompanionUreadCount}, nil
}

func (c *Core) GetPrivateRoomWithoutCreate(ctx context.Context, user1 uint64, user2 uint64) (*entity.Room, error) {
	if user1 == user2 {
		return nil, entity.CreateError(400, entity.ErrRefIDEquals)
	}
	return c.Storage.GetPrivateRoom(ctx, user1, user2)
}

func (c *Core) SetZeroUnreadCount(ctx context.Context, roomid string, who uint64) error {
	if err := c.Storage.ZeroUnreadCountInRoom(ctx, roomid, who); err != nil {
		return err
	}
	return c.sendMyRoom(ctx, roomid)
}
