package configservice

import (
	"log"
	"os"

	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
	"github.com/irrepressiblespirit/website-messages-service/pkg/service"
	"gopkg.in/yaml.v2"
)

type ConfigService struct {
	FileName string
}

func NewConfigService(fileName string) service.ConfigService {
	return ConfigService{
		FileName: fileName,
	}
}

func (service ConfigService) Load() (*entity.Config, error) {
	log.Printf("Loading config from %s", service.FileName)
	config := &entity.Config{}
	file, err := os.Open(service.FileName)
	if err != nil {
		return config, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return config, err
	}
	return config, err
}
