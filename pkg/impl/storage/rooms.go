package storage

import (
	"context"
	"errors"
	"time"

	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *MongoStorage) GetPrivateRoom(ctx context.Context, user1 uint64, user2 uint64) (*entity.Room, error) {
	// example: db.getCollection('rooms').find({"users.refid": {$all: [2525,4546]}})
	room := entity.Room{}
	err := s.database.Collection(MongoCollectionRooms).FindOne(
		ctx,
		bson.M{"users.refid": bson.M{"$all": []uint64{user1, user2}}},
	).Decode(&room)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.RoomNotFoundError{}
		}
		return nil, err
	}
	return &room, nil
}

func (s *MongoStorage) CreatePrivateRoom(ctx context.Context, room *entity.Room) (string, error) {
	res, err := s.database.Collection(MongoCollectionRooms).InsertOne(ctx, room)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *MongoStorage) GetRoom(ctx context.Context, roomid string) (*entity.Room, error) {
	room := entity.Room{}
	rid, err := primitive.ObjectIDFromHex(roomid)
	if err != nil {
		return nil, err
	}
	err = s.database.Collection(MongoCollectionRooms).FindOne(
		ctx, bson.M{"_id": rid}).Decode(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *MongoStorage) GetUserRooms(ctx context.Context, user uint64, page int64, size int64) ([]*entity.Room, error) {
	// db.getCollection('rooms').find({"users.refid": 13706951275})
	// db.getCollection('rooms').find({"users.refid": 13706951275}).skip(1).limit(1)
	res := []*entity.Room{}
	options := options.Find()
	options.SetSort(bson.M{"last_msg": -1})
	options.SetLimit(size)
	options.SetSkip(page * size)
	cur, err := s.database.Collection(MongoCollectionRooms).Find(
		ctx,
		bson.M{"users.refid": user},
		options,
	)
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var room entity.Room
		err = cur.Decode(&room)
		if err != nil {
			return nil, err
		}
		res = append(res, &room)
	}
	return res, nil
}

func (s *MongoStorage) IncreaseUnreadCountInRoom(ctx context.Context, roomid string) error {
	/*
		db.getCollection('rooms').update(
			{"_id": ObjectId("1111")},
			{ $inc: {"users.$[].notread": 1} }
		)
	*/
	rid, err := primitive.ObjectIDFromHex(roomid)
	if err != nil {
		return err
	}
	_, err = s.database.Collection(MongoCollectionRooms).UpdateMany(
		ctx,
		bson.M{"_id": rid},
		bson.M{"$inc": bson.M{"users.$[].notread": 1}},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoStorage) GetUnreadCountMessagesByAllRooms(ctx context.Context, user uint64) (int64, error) {
	// db.getCollection('rooms').aggregate([{$unwind: "$users"},{ $match: {"users.refid": 1597}}, { $group: { "_id" : null, "sum" : { $sum: "$users.notread" } }}])
	var results []struct {
		ID       string `bson:"_id"`
		TotalSum int64  `bson:"sum"`
	}
	unwindStage := bson.M{"$unwind": "$users"}
	matchStage := bson.M{"$match": bson.M{"users.refid": user}}
	groupStage := bson.M{"$group": bson.M{"_id": nil, "sum": bson.M{"$sum": "$users.notread"}}}
	c, err := s.database.Collection(MongoCollectionRooms).Aggregate(
		ctx,
		[]bson.M{unwindStage, matchStage, groupStage},
	)
	if err != nil {
		return int64(0), err
	}
	if err = c.All(ctx, &results); err != nil {
		return int64(0), err
	}
	if len(results) > 0 {
		return results[0].TotalSum, nil
	}
	return int64(0), nil
}

func (s *MongoStorage) SetCurrentTimeInRoom(ctx context.Context, roomid string) error {
	/*
		db.getCollection('rooms').update(
			{"_id": ObjectId("2222")},
			{$set: {"last_msg": 111111}}
		)
	*/
	rid, err := primitive.ObjectIDFromHex(roomid)
	if err != nil {
		return err
	}
	_, err = s.database.Collection(MongoCollectionRooms).UpdateOne(
		ctx,
		bson.M{"_id": rid},
		bson.M{"$set": bson.M{"last_msg": primitive.Timestamp{T: uint32(time.Now().Unix())}}},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoStorage) ZeroUnreadCountInRoom(ctx context.Context, roomid string, refid uint64) error {
	/*
	   db.getCollection('rooms').update(
	     {"_id": ObjectId("1111"), "users.refid": 2324},
	     { $set: { "users.$.notread": 0 }}
	   )
	*/
	rid, err := primitive.ObjectIDFromHex(roomid)
	if err != nil {
		return err
	}

	_, err = s.database.Collection(MongoCollectionRooms).UpdateMany(
		ctx,
		bson.M{"_id": rid, "users.refid": refid},
		bson.M{"$set": bson.M{"users.$.notread": 0}},
	)
	if err != nil {
		return err
	}
	return nil
}
