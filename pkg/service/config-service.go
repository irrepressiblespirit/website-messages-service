package service

import "github.com/irrepressiblespirit/website-messages-service/pkg/entity"

type ConfigService interface {
	Load() (*entity.Config, error)
}
