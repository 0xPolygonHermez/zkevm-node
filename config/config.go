package config

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/aggregator"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/ethermanv2"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/gasprice"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpcv2"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter"
	"github.com/0xPolygonHermez/zkevm-node/proverclient"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/sequencerv2"
	"github.com/0xPolygonHermez/zkevm-node/sequencerv2/broadcast"
	"github.com/0xPolygonHermez/zkevm-node/state/tree"
	"github.com/0xPolygonHermez/zkevm-node/statev2/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer"
	"github.com/0xPolygonHermez/zkevm-node/synchronizerv2"
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
	Ethermanv2        ethermanv2.Config
	EthTxManager      ethtxmanager.Config
	RPC               jsonrpc.Config
	RPCV2             jsonrpcv2.Config
	Synchronizer      synchronizer.Config
	Synchronizerv2    synchronizerv2.Config
	Sequencer         sequencer.Config
	SequencerV2       sequencerv2.Config
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
