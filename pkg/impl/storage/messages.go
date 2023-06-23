package storage

import (
	"context"

	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *MongoStorage) SaveMsg(ctx context.Context, msg entity.IMessage) (string, error) {
	res, err := s.database.Collection(MongoCollectionMessages).InsertOne(ctx, msg)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *MongoStorage) GetMsgs(ctx context.Context, roomid string, limit int64, lastMsgID string) ([]entity.IMessage, error) {
	// db.getCollection('messages').find({"roomid": "1111"})
	var findFilter primitive.M
	if lastMsgID == "" {
		findFilter = bson.M{"roomid": roomid}
	} else {
		lastMsgIDObjectID, err := primitive.ObjectIDFromHex(lastMsgID)
		if err != nil {
			return []entity.IMessage{}, err
		}
		findFilter = bson.M{"roomid": roomid, "_id": bson.M{"$lt": lastMsgIDObjectID}}
	}
	findOpts := options.Find()
	findOpts.SetLimit(limit)
	findOpts.SetSort(bson.M{"sended": -1})
	res := []entity.IMessage{}
	cur, err := s.database.Collection(MongoCollectionMessages).Find(
		ctx,
		findFilter,
		findOpts,
	)
	if err != nil {
		return []entity.IMessage{}, err
	}
	for cur.Next(ctx) {
		var tmp entity.Message
		err := cur.Decode(&tmp)
		if err != nil {
			return []entity.IMessage{}, err
		}
		msg, err := entity.GetMessageStructByType(tmp.Type)
		if err != nil {
			return []entity.IMessage{}, err
		}
		err = cur.Decode(msg)
		if err != nil {
			return []entity.IMessage{}, err
		}
		res = append(res, msg)
	}
	return res, nil
}

func (s *MongoStorage) GetLastMsg(ctx context.Context, roomid string) (entity.IMessage, error) {
	options := options.FindOne()
	options.SetSort(bson.M{"sended": -1})
	var tmp entity.Message
	sr := s.database.Collection(MongoCollectionMessages).FindOne(
		ctx,
		bson.M{"roomid": roomid},
		options,
	)
	err := sr.Decode(&tmp)
	if err != nil {
		return nil, err
	}
	msg, err := entity.GetMessageStructByType(tmp.GetType())
	if err != nil {
		return nil, err
	}
	err = sr.Decode(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
