package grpcapi

import (
	"context"
	"encoding/json"

	"github.com/irrepressiblespirit/website-messages-service/pkg/core"
	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
	pb "github.com/irrepressiblespirit/website-messages-service/pkg/grpcapi"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	Core core.Core
	pb.UnimplementedServiceMessagesServer
}

func (s *GRPCServer) RequestNewToken(ctx context.Context, in *pb.TokenRequest) (*pb.TokenResponse, error) {
	token, err := s.Core.GetNewToken(uint64(in.GetRefid()))
	if err != nil {
		return nil, parseError(err)
	}
	centrifugoURL, err := entity.GetCentrifugoURL(s.Core.Config)
	if err != nil {
		return nil, parseError(err)
	}
	return &pb.TokenResponse{
		Refid:         int64(token.RefID),
		Token:         token.Token,
		CentrifugoURL: centrifugoURL,
	}, nil
}

func (s *GRPCServer) GetMyRooms(ctx context.Context, in *pb.MyRoomsRequest) (*pb.MyRoomsResponse, error) {
	rooms, count, err := s.Core.GetMyRooms(ctx, uint64(in.GetRefid()), in.Page, in.Size)
	if err != nil {
		return nil, parseError(err)
	}
	result := pb.MyRoomsResponse{}
	for i := 0; i < len(rooms); i++ {
		tmp := pb.MyRoom{
			RoomID:              rooms[i].ID,
			Type:                rooms[i].Type,
			Name:                rooms[i].Name,
			LastMessage:         timestamppb.New(rooms[i].LastMessage),
			LogoUrl:             rooms[i].LogoURL,
			NotReadCount:        int32(rooms[i].NotReadCount),
			CompanionRefID:      int64(rooms[i].CompanionRefID),
			CompanionOwnerRefID: int64(rooms[i].CompanionOwnerRefID),
		}
		if rooms[i].LastMsg != nil {
			tmp.MyRoomLastMessage = &pb.MyRoomLastMessage{
				ID:    rooms[i].LastMsg.ID,
				RefID: int64(rooms[i].LastMsg.RefID),
				Type:  string(rooms[i].LastMsg.Type),
				Body:  rooms[i].LastMsg.Body,
			}
		}
		result.Rooms = append(result.Rooms, &tmp)
	}

	result.UnreadMessagesCount = count
	return &result, nil
}

func (s *GRPCServer) SendFastMessage(ctx context.Context, in *pb.FastMessageRequest) (*pb.FastMessageResponse, error) {
	roomID, msgID, err := s.Core.SendFastMsg(
		ctx,
		uint64(in.GetRefIDFrom()),
		uint64(in.GetRefIDTo()),
		in.GetMessageBody(),
	)
	if err != nil {
		return nil, parseError(err)
	}
	return &pb.FastMessageResponse{ID: msgID, RoomID: roomID}, nil
}

func (s *GRPCServer) SendMessageInRoom(ctx context.Context, in *pb.SendMessageInRoomRequest) (*pb.SendMessageInRoomResponse, error) {
	msg, err := entity.GetMessageStructByType(entity.MessageType(in.MessageType))
	if err != nil {
		return nil, parseError(err)
	}
	msg.SetVersion(entity.CurrentMessageVersion).
		SetRefID(uint64(in.GetRefIDFrom())).
		SetRoomID(in.RoomID).
		SetType(entity.MessageType(in.MessageType)).
		SetBody(in.MessageBody)
	b, err := json.Marshal(in.Addition)
	if err != nil {
		return nil, parseError(err)
	}
	addMap := make(map[string]interface{})
	err = json.Unmarshal(b, &addMap)
	if err != nil {
		return nil, parseError(err)
	}
	if err := msg.ParseAddition(addMap); err != nil {
		return nil, parseError(err)
	}
	msgID, err := s.Core.SendMsg(ctx, msg)
	if err != nil {
		return nil, parseError(err)
	}
	return &pb.SendMessageInRoomResponse{ID: msgID}, nil
}

func (s *GRPCServer) SendMessageByRefID(ctx context.Context, in *pb.SendMessageByRefIDRequest) (*pb.SendMessageByRefIDResponse, error) {
	msg, err := entity.GetMessageStructByType(entity.MessageType(in.MessageType))
	if err != nil {
		return nil, parseError(err)
	}
	room, err := s.Core.GetPrivateRoom(ctx, uint64(in.GetRefIDFrom()), uint64(in.GetRefIDTo()))
	if err != nil {
		return nil, parseError(err)
	}
	msg.SetVersion(entity.CurrentMessageVersion).
		SetRefID(uint64(in.GetRefIDFrom())).
		SetRoomID(room.ID).
		SetType(entity.MessageType(in.MessageType)).
		SetBody(in.MessageBody)
	b, err := json.Marshal(in.Addition)
	if err != nil {
		return nil, parseError(err)
	}
	addMap := make(map[string]interface{})
	err = json.Unmarshal(b, &addMap)
	if err != nil {
		return nil, parseError(err)
	}
	if err := msg.ParseAddition(addMap); err != nil {
		return nil, parseError(err)
	}
	msgID, err := s.Core.SendMsg(ctx, msg)
	if err != nil {
		return nil, parseError(err)
	}
	return &pb.SendMessageByRefIDResponse{ID: msgID, RoomID: room.ID}, nil
}

func (s *GRPCServer) GetMessagesInRoom(ctx context.Context, in *pb.MessagesRequest) (*pb.MessagesResponse, error) {
	messages, err := s.Core.GetMessages(
		ctx,
		in.GetRoomID(),
		uint64(in.GetRefID()),
		int(in.GetCount()),
		in.GetLastMsgID(),
	)
	if err != nil {
		return nil, parseError(err)
	}
	result := pb.MessagesResponse{}
	result.CompanionUreadCount = int32(messages.CompanionUreadCount)
	for i := 0; i < len(messages.Messages); i++ {
		msg := pb.Message{
			ID:          messages.Messages[i].GetID(),
			RoomID:      messages.Messages[i].GetRoomID(),
			RefID:       int64(messages.Messages[i].GetRefID()),
			Sended:      timestamppb.New(messages.Messages[i].GetSended()),
			MessageType: string(messages.Messages[i].GetType()),
			Body:        messages.Messages[i].GetBody(),
		}
		b, err := json.Marshal(messages.Messages[i].GetAddition())
		if err != nil {
			return nil, parseError(err)
		}
		msg.Addition = &structpb.Struct{}
		err = json.Unmarshal(b, &msg.Addition.Fields)
		if err != nil {
			return nil, parseError(err)
		}
		result.Messages = append(result.Messages, &msg)
	}
	return &result, nil
}

func (s *GRPCServer) GetUserInfo(ctx context.Context, in *pb.UserRequest) (*pb.User, error) {
	user, err := s.Core.GetUser(ctx, uint64(in.GetRefID()))
	if err != nil {
		return nil, parseError(err)
	}
	return &pb.User{
		RefID:      int64(user.RefID),
		Name:       user.Name,
		LogoUrl:    user.LogoURL,
		CachedTime: timestamppb.New(user.CachedTime),
		OwnerRefID: int64(user.OwnerRefID),
	}, nil
}

func (s *GRPCServer) GetPrivateRoomWithoutCreate(ctx context.Context, in *pb.GetPrivateRoomWithoutCreateRequest) (*pb.MyRoom, error) {
	room, err := s.Core.GetPrivateRoomWithoutCreate(ctx, uint64(in.GetRefID()), uint64(in.GetRefID2()))
	if err != nil {
		return nil, parseError(err)
	}
	myroom, err := s.Core.ConvertPrivateRoomToMyRoom(ctx, room, uint64(in.GetRefID()))
	if err != nil {
		return nil, parseError(err)
	}
	return &pb.MyRoom{
		RoomID:              myroom.ID,
		Type:                myroom.Type,
		Name:                myroom.Name,
		LastMessage:         timestamppb.New(myroom.LastMessage),
		LogoUrl:             myroom.LogoURL,
		NotReadCount:        int32(myroom.NotReadCount),
		CompanionRefID:      int64(myroom.CompanionRefID),
		CompanionOwnerRefID: int64(myroom.CompanionOwnerRefID),
		MyRoomLastMessage: &pb.MyRoomLastMessage{
			ID:    myroom.LastMsg.ID,
			RefID: int64(myroom.LastMsg.RefID),
			Type:  string(myroom.LastMsg.Type),
			Body:  myroom.LastMsg.Body,
		},
	}, nil
}

func (s *GRPCServer) SetZeroUnreadCount(ctx context.Context, in *pb.SetZeroUnreadCountRequest) (*pb.SetZeroUnreadCountResponse, error) {
	err := s.Core.SetZeroUnreadCount(ctx, in.RoomID, uint64(in.RefID))
	if err != nil {
		return nil, parseError(err)
	}
	return &pb.SetZeroUnreadCountResponse{}, nil
}
