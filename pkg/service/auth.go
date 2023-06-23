package service

type IAuthToken interface {
	Check(token string) error
}
