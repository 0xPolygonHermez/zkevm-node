package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/hermeznetwork/hermez-core/aggregator"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	logger "github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/synchronizer"
	configLibrary "github.com/hermeznetwork/go-hermez-config"
	"github.com/go-playground/validator/v10"
)

// Config represents the configuration of the entire Hermez Node
type Config struct {
	Log          logger.Config
	Database     db.Config
	Synchronizer synchronizer.Config
	RPC          jsonrpc.Config
	Sequencer    sequencer.Config
	Aggregator   aggregator.Config
	Prover       proverclient.Config
}

// Load loads the configuration
func Load(configFilePath string) (*Config, error) {
	var cfg Config
	err := configLibrary.LoadConfig(configFilePath, DefaultValues, &cfg)
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
	log.Printf("Configuration loaded: %+v", cfg)
	return &cfg, nil
}
