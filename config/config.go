package config

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/hermeznetwork/hermez-core/aggregator"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/gasprice"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/synchronizer"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	flagCfg        = "cfg"
	flagNetwork    = "network"
	flagNetworkCfg = "network-cfg"
)

// Config represents the configuration of the entire Hermez Node
type Config struct {
	Log               log.Config
	Database          db.Config
	Etherman          etherman.Config
	RPC               jsonrpc.Config
	Synchronizer      synchronizer.Config
	Sequencer         sequencer.Config
	Aggregator        aggregator.Config
	Prover            proverclient.Config
	NetworkConfig     NetworkConfig
	GasPriceEstimator gasprice.Config
	MTService         tree.Config
}

// Load loads the configuration
func Load(ctx *cli.Context) (*Config, error) {
	var cfg Config
	viper.SetConfigType("toml")

	err := viper.ReadConfig(bytes.NewBuffer([]byte(DefaultValues)))
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, err
	}
	configFilePath := ctx.String(flagCfg)
	if configFilePath != "" {
		dirName, fileName := filepath.Split(configFilePath)

		fileExtension := strings.TrimPrefix(filepath.Ext(fileName), ".")
		fileNameWithoutExtension := strings.TrimSuffix(fileName, "."+fileExtension)

		viper.AddConfigPath(dirName)
		viper.SetConfigName(fileNameWithoutExtension)
		viper.SetConfigType(fileExtension)
	}
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("HERMEZCORE")
	err = viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			log.Infof("config file not found")
		} else {
			log.Infof("error reading config file: ", err)
			return nil, err
		}
	}

	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, err
	}
	// Load genesis parameters
	cfg.loadNetworkConfig(ctx)

	cfgJSON, _ := json.MarshalIndent(cfg, "", "  ")
	log.Infof("Configuration loaded: \n%s\n", string(cfgJSON))
	return &cfg, nil
}
