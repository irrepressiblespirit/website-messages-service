package tokenservice

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
	"github.com/irrepressiblespirit/website-messages-service/pkg/service"
	"github.com/irrepressiblespirit/website-messages-service/pkg/utils"
)

const MaxJwtTokenTime = 15

type tokenServiceImpl struct {
	secret []byte
}

func NewTokenService(secret string) service.TokenService {
	return tokenServiceImpl{
		secret: []byte(secret),
	}
}

func (t tokenServiceImpl) Generate(refid uint64) (entity.Token, error) {
	return t.GenerateWithTime(refid, time.Now().Add(time.Minute*MaxJwtTokenTime).Unix())
}

func (t tokenServiceImpl) GenerateWithTime(refid uint64, timeExp int64) (entity.Token, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": utils.Uint64ToString(refid),
			"exp": timeExp,
		})
	tokenString, err := token.SignedString(t.secret)
	if err != nil {
		return entity.Token{}, entity.CoreError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return entity.Token{
		RefID: refid,
		Token: tokenString,
	}, nil
}
