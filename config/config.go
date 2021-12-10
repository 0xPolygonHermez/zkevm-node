package config

import (
	"bytes"
	"log"
	"path/filepath"
	"strings"

	"github.com/hermeznetwork/hermez-core/aggregator"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	logger "github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/synchronizer"
	"github.com/spf13/viper"
	"github.com/mitchellh/mapstructure"
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
	viper.SetConfigType("toml")

	viper.ReadConfig(bytes.NewBuffer([]byte(DefaultValues)))
	err := viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, err
	}
	path, fullFile := filepath.Split(configFilePath)

	file := strings.Split(fullFile, ".")

	viper.AddConfigPath(path)
	viper.SetConfigName(file[0])
	viper.SetConfigType(file[1])
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("HERMEZCORE")
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, err
	}
	log.Printf("Configuration loaded: %+v", cfg)
	return &cfg, nil
}
