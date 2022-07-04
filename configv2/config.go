package config

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"

	aggregator "github.com/hermeznetwork/hermez-core/aggregatorv2"
	"github.com/hermeznetwork/hermez-core/db"
	etherman "github.com/hermeznetwork/hermez-core/ethermanv2"
	"github.com/hermeznetwork/hermez-core/ethtxmanager"
	gasprice "github.com/hermeznetwork/hermez-core/gaspricev2"
	jsonrpc "github.com/hermeznetwork/hermez-core/jsonrpcv2"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pricegetter"
	"github.com/hermeznetwork/hermez-core/proverclient"
	sequencer "github.com/hermeznetwork/hermez-core/sequencerv2"
	"github.com/hermeznetwork/hermez-core/sequencerv2/broadcast"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/executor"
	synchronizer "github.com/hermeznetwork/hermez-core/synchronizerv2"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	FlagYes         = "yes"
	FlagCfg         = "cfg"
	FlagNetwork     = "network"
	FlagNetworkCfg  = "network-cfg"
	FlagNetworkBase = "network-base"
	FlagAmount      = "amount"
	FlagRemoteMT    = "remote-merkletree"
	FlagComponents  = "components"
	FlagHTTPAPI     = "http.api"
)

// Config represents the configuration of the entire Hermez Node
type Config struct {
	Log               log.Config
	Database          db.Config
	Etherman          etherman.Config
	EthTxManager      ethtxmanager.Config
	RPC               jsonrpc.Config
	Synchronizer      synchronizer.Config
	Sequencer         sequencer.Config
	PriceGetter       pricegetter.Config
	Aggregator        aggregator.Config
	Prover            proverclient.Config
	NetworkConfig     NetworkConfig
	GasPriceEstimator gasprice.Config
	MTServer          tree.ServerConfig
	MTClient          tree.ClientConfig
	Executor          executor.Config
	BroadcastServer   broadcast.ServerConfig
	BroadcastClient   broadcast.ClientConfig
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
	configFilePath := ctx.String(FlagCfg)
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
