package authtoken

import (
	"net/http"

	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
)

type AuthToken struct {
	Token string
}

func New(token string) *AuthToken {
	return &AuthToken{
		Token: token,
	}
}

func (a *AuthToken) Check(token string) error {
	if a.Token == "" {
		return nil
	}
	if a.Token == token {
		return nil
	}
	return entity.CoreError{Code: http.StatusInternalServerError, Message: entity.ErrInvalidToken}
}
