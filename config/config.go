package config

import (
	"log"
	"strings"
	configLibrary "github.com/hermeznetwork/go-hermez-config"
	"gopkg.in/go-playground/validator.v9"
)

type Conf struct {}

// LoadConfig loads the configuration from path.
func LoadConfig(path string) (*Conf, error) {
	var cfg Conf
	err := configLibrary.LoadConfig(path, DefaultValues, &cfg)
	if err != nil {
		//Split errors depending on if there is a file error, a env error or a default error
		if strings.Contains(err.Error(), "default") {
			return nil, err
		}
		log.Println(err.Error())
	}
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("error validating configuration file: %w", err)
	}
	log.Printf("Loaded Configuration: %+v", cfg)
	return &cfg, nil
}