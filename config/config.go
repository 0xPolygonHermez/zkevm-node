package config

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/aggregator"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/gasprice"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter"
	"github.com/0xPolygonHermez/zkevm-node/proverclient"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	// FlagYes is the flag for yes.
	FlagYes = "yes"
	// FlagCfg is the flag for cfg.
	FlagCfg = "cfg"
	// FlagGenesisFile is the flag for genesis file.
	FlagGenesisFile = "genesis"
	// FlagAmount is the flag for amount.
	FlagAmount = "amount"
	// FlagRemoteMT is the flag for remote-merkletree.
	FlagRemoteMT = "remote-merkletree"
	// FlagComponents is the flag for components.
	FlagComponents = "components"
	// FlagHTTPAPI is the flag for http.api.
	FlagHTTPAPI = "http.api"
)

// Config represents the configuration of the entire Hermez Node
type Config struct {
	IsTrustedSequencer bool `mapstructure:"IsTrustedSequencer"`
	Log                log.Config
	Etherman           etherman.Config
	EthTxManager       ethtxmanager.Config
	RPC                jsonrpc.Config
	Synchronizer       synchronizer.Config
	Sequencer          sequencer.Config
	PriceGetter        pricegetter.Config
	Aggregator         aggregator.Config
	Provers            proverclient.Config
	NetworkConfig      NetworkConfig
	GasPriceEstimator  gasprice.Config
	Executor           executor.Config
	BroadcastServer    broadcast.ServerConfig
	MTClient           merkletree.Config
	StateDB            db.Config
	PoolDB             db.Config
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
	viper.SetEnvPrefix("ZKEVM_NODE")
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
	/*
		cfgJSON, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return nil, err
		}
		log.Debugf("Configuration loaded: \n%s\n", string(cfgJSON))
	*/
	return &cfg, nil
}
