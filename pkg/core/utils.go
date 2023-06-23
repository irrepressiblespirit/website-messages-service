package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
)

func (c *Core) ConvertPrivateRoomToMyRoom(ctx context.Context, room *entity.Room, whosee uint64) (*entity.MyRoom, error) {
	myroom := entity.MyRoom{
		ID:           room.ID,
		Type:         room.Type,
		LastMessage:  room.LastMessage,
		Name:         "not ready",
		LogoURL:      "not ready",
		NotReadCount: 0,
	}
	if me := room.GetUser(whosee); me != nil {
		myroom.NotReadCount = me.NotReadCount
	}
	other := room.GetOtherUser(whosee)
	if other != nil {
		otherUser, err := c.GetUser(ctx, other.RefID)
		if err != nil {
			var errUNF entity.UserNotFoundError
			if errors.Is(err, &errUNF) {
				otherUser = entity.GetUserNotFoundItem(other.RefID)
			} else {
				return nil, err
			}
		}
		myroom.Name = otherUser.Name
		myroom.LogoURL = otherUser.LogoURL
		myroom.CompanionRefID = otherUser.RefID
		myroom.CompanionOwnerRefID = otherUser.OwnerRefID
		myroom.CompanionUnreadCount = other.NotReadCount
	}
	msg, err := c.Storage.GetLastMsg(ctx, room.ID)
	if msg == nil {
		myroom.LastMsg = nil
		return &myroom, nil
	}
	if err != nil {
		return nil, err
	}
	myroom.LastMsg = &entity.MyRoomLastMsg{
		ID:    msg.GetID(),
		RefID: msg.GetRefID(),
		Type:  msg.GetType(),
		Body:  msg.GetBody(),
	}
	return &myroom, nil
}

func (c *Core) sendMessageInRoom(ctx context.Context, room *entity.Room, msg entity.IMessage) error {
	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for i := 0; i < len(room.Users); i++ {
		err := c.EventStorage.SendMessage(ctx, fmt.Sprintf("messages#%d", room.Users[i].RefID), json)
		if err != nil {
			return err
		}
		// fmt.Printf("Event [messages#%d] -> %s\n", room.Users[i].RefID, string(json)) //TODO в centrifugo
	}
	return nil
}

func (c *Core) sendMyRoom(ctx context.Context, roomid string) error {
	room, err := c.Storage.GetRoom(ctx, roomid)
	if err != nil {
		return err
	}
	switch room.Type {
	case entity.RoomTypePrivate:
		for _, user := range room.Users {
			myroom, err := c.ConvertPrivateRoomToMyRoom(ctx, room, user.RefID)
			if err != nil {
				return err
			}
			count, _ := c.GetUnreadCount(ctx, user.RefID)
			body := entity.MyRoomWithNotReadMessagesTotalCount{
				Room:                 *myroom,
				NotReadMessagesCount: count,
			}
			json, err2 := json.Marshal(&body)
			if err2 != nil {
				return err2
			}
			err = c.EventStorage.SendMessage(ctx, fmt.Sprintf("rooms#%d", user.RefID), json)
			if err != nil {
				return err
			}
			// fmt.Printf("Event [rooms#%d] -> %s\n", user.RefID, string(json)) //TODO в центрифугу
		}
	default:
		return entity.CreateError(500, entity.ErrRoomsUnsupportRoomType)
	}
	return nil
}
