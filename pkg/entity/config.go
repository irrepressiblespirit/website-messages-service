package entity

import (
	"crypto/rand"
	"math/big"
	"time"
)

type Config struct {
	Prefix       string           `yaml:"prefix"`
	Token        string           `yaml:"token"`
	Mongo        string           `yaml:"mongo"`
	Centrifugo   ConfigCentrifugo `yaml:"centrifugo"`
	ExternalUser string           `yaml:"externaluser"`
	ConfigGin    ConfigGin        `yaml:"gin"`
	RestTimeout  time.Duration    `yaml:"resttimeout"`
	WS           []string         `yaml:"ws"`
}

type ConfigCentrifugo struct {
	Secret string `yaml:"secret"`
	APIKEY string `yaml:"apikey"`
	APIURL string `yaml:"apiurl"`
	Grpc   string `yaml:"grpc"`
}

type ConfigGin struct {
	Mode string `yaml:"mode"`
}

func GetCentrifugoURL(conf *Config) (string, error) {
	limit := int64(len(conf.WS))
	randomValue, err := rand.Int(rand.Reader, big.NewInt(limit))
	if err != nil {
		return "", err
	}
	return conf.WS[randomValue.Int64()], nil
}
