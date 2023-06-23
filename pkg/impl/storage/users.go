package storage

import (
	"context"
	"errors"
	"time"

	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
	"github.com/ulule/deepcopier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type storedUser struct {
	RefID      uint64    `json:"refid" deepcopier:"field:RefID"`
	Name       string    `json:"name" deepcopier:"field:Name"`
	LogoURL    string    `json:"logourl" deepcopier:"field:LogoUrl"`
	CachedTime time.Time `json:"cached_time" deepcopier:"field:CachedTime"`
	OwnerRefID uint64    `json:"ownerrefid" deepcopier:"field:OwnerRefId"`
}

func (s *MongoStorage) GetUser(ctx context.Context, refid uint64) (*entity.User, error) {
	var storedUser storedUser
	var user entity.User

	err := s.database.Collection(MongoCollectionUsers).
		FindOne(ctx, bson.M{"refid": refid}).Decode(&storedUser)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &entity.UserNotFoundError{}
		}
		return nil, err
	}
	err = deepcopier.Copy(&storedUser).To(&user)
	if err != nil {
		return nil, err
	}
	user.System.Cached = true
	return &user, nil
}

func (s *MongoStorage) PutUser(ctx context.Context, user *entity.User) error {
	storedUser := storedUser{}
	err := deepcopier.Copy(user).To(&storedUser)
	if err != nil {
		return err
	}
	_, err = s.database.Collection(MongoCollectionUsers).InsertOne(ctx, &storedUser)
	if err != nil {
		return err
	}
	return nil
}
