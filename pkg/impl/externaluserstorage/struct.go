package externaluserstorage

import (
	"strconv"
	"time"

	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
)

func (extUser ExternalServiceUser) ConvertToUser() *entity.User {
	refId, _ := strconv.ParseUint(extUser.Results[0].Id.Value, 10, 64)
	return &entity.User{
		RefID:      refId,
		Name:       extUser.Results[0].Name.First + " " + extUser.Results[0].Name.Last,
		LogoURL:    extUser.Results[0].Picture.Medium,
		CachedTime: time.Now(),
		OwnerRefID: 0,
	}
}

type ExternalServiceUser struct {
	Results []externalUserItem `json:"results"`
	Info    infoStruct         `json:"info"`
}

type externalUserItem struct {
	Gender string `json:"gender"`
	Name   struct {
		Title string `json:"title"`
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
	Location struct {
		Street struct {
			Number int    `json:"number"`
			Name   string `json:"name"`
		} `json:"street"`
		City        string `json:"city"`
		State       string `json:"state"`
		Country     string `json:"country"`
		Postcode    int    `json:"postcode"`
		Coordinates struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"coordinates"`
		Timezone struct {
			Offset      string `json:"offset"`
			Description string `json:"description"`
		} `json:"timezone"`
	} `json:"location"`
	Email string `json:"email"`
	Login struct {
		Uuid     string `json:"uuid"`
		Username string `json:"username"`
		Password string `json:"password"`
		Salt     string `json:"salt"`
		Md5      string `json:"md5"`
		Sha1     string `json:"sha1"`
		Sha256   string `json:"sha256"`
	} `json:"login"`
	Dob struct {
		Date string `json:"date"`
		Age  int    `json:"age"`
	} `json:"dob"`
	Registered struct {
		Date string `json:"date"`
		Age  int    `json:"age"`
	} `json:"registered"`
	Phone string `json:"phone"`
	Cell  string `json:"cell"`
	Id    struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"id"`
	Picture struct {
		Large     string `json:"large"`
		Medium    string `json:"medium"`
		Thumbnail string `json:"thumbnail"`
	} `json:"picture"`
	Nat string `json:"nat"`
}

type infoStruct struct {
	Seed    string `json:"seed"`
	Results int    `json:"results"`
	Page    int    `json:"page"`
	Version string `json:"version"`
}
