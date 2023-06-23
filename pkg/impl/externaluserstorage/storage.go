package externaluserstorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
	"github.com/irrepressiblespirit/website-messages-service/pkg/service"
)

const execTimeOut = time.Second * 3

type StorageUserDecorator struct {
	URL     string
	storage service.StorageUsers
	rest    *resty.Client
}

func NewStorageUserDecorator(url string, storage service.StorageUsers) *StorageUserDecorator {
	return &StorageUserDecorator{
		URL:     url,
		storage: storage,
		rest:    resty.New().SetTimeout(execTimeOut),
	}
}

func (s *StorageUserDecorator) GetUser(ctx context.Context, refid uint64) (*entity.User, error) {
	user, err := s.storage.GetUser(ctx, refid)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, &entity.UserNotFoundError{}) {
		user, err := s.getUserFromExternalService(ctx, refid)
		if err != nil {
			return nil, fmt.Errorf("error StorageUserDecorator.GetUser: %w [%d]", err, refid)
		}
		if err = s.storage.PutUser(ctx, user); err != nil {
			return nil, fmt.Errorf("error StorageUserDecorator.GetUser: %w [%d]", err, refid)
		}
		return user, nil
	}
	return nil, fmt.Errorf("error StorageUserDecorator.GetUser: %w", err)
}

func (s *StorageUserDecorator) PutUser(ctx context.Context, user *entity.User) error {
	return s.storage.PutUser(ctx, user)
}

func (s StorageUserDecorator) getUserFromExternalService(ctx context.Context, refid uint64) (*entity.User, error) {
	// add code which get user info from external API by user refId logic
	// now in code i send a request to random user generator (https://randomuser.me/api/), get user and convert him
	resp, err := s.rest.R().
		SetHeader("Content-Type", "text/plain;charset=UTF-8").
		SetContext(ctx).
		Get(s.URL)
	if err != nil {
		return nil, err
	}
	var externalUser ExternalServiceUser
	err = json.Unmarshal(resp.Body(), &externalUser)
	if err != nil {
		return nil, err
	}
	return externalUser.ConvertToUser(), nil
}
