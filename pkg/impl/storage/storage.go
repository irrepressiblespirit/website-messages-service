package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoDBDatabase         = "messages"
	MongoCollectionRooms    = "rooms"
	MongoCollectionMessages = "messages"
	MongoCollectionUsers    = "users"

	timeoutConnect = time.Second * 1
)

type MongoStorage struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoStorage(uri string) (*MongoStorage, error) {
	storage := MongoStorage{}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutConnect)
	defer cancel()
	client, err := mongo.NewClient(
		options.
			Client().
			ApplyURI(uri).
			SetConnectTimeout(1 * time.Second),
	)
	if err != nil {
		return &MongoStorage{}, err
	}
	err = client.Connect(ctx)
	if err != nil {
		return &MongoStorage{}, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return &MongoStorage{}, err
	}
	storage.client = client
	storage.database = client.Database(MongoDBDatabase)
	return &storage, nil
}

func (s *MongoStorage) Ping(ctx context.Context) error {
	return s.client.Ping(ctx, nil)
}
