package entity

import "time"

type User struct {
	RefID      uint64     `json:"refid" deepcopier:"field:RefID"`
	Name       string     `json:"name" deepcopier:"field:Name"`
	LogoURL    string     `json:"logourl" deepcopier:"field:LogoUrl"`
	CachedTime time.Time  `json:"cached_time" deepcopier:"field:CachedTime"`
	System     UserSystem `json:"system"`
	OwnerRefID uint64     `json:"ownerrefid" deepcopier:"field:OwnerRefId"`
}

type UserSystem struct {
	Cached bool `json:"cached" deepcopier:"skip"`
}

type UserNotFoundError struct{}

func (e *UserNotFoundError) Error() string {
	return "User not found"
}

func GetUserNotFoundItem(refID uint64) *User {
	return &User{
		RefID:      refID,
		Name:       "name.user.not.found",
		LogoURL:    "",
		CachedTime: time.Now(),
		OwnerRefID: refID,
	}
}
