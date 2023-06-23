package service

import "github.com/irrepressiblespirit/website-messages-service/pkg/entity"

type TokenService interface {
	Generate(refid uint64) (entity.Token, error)
	GenerateWithTime(refid uint64, timeExp int64) (entity.Token, error)
}
