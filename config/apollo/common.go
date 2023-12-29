package apollo

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"os"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func (c *Client) unmarshal(value interface{}) (*config.Config, error) {
	v := viper.New()
	v.SetConfigType("toml")
	err := v.ReadConfig(bytes.NewBuffer([]byte(value.(string))))
	if err != nil {
		log.Errorf("failed to load config: %v error: %v", value, err)
		return nil, err
	}
	dstConf := config.Config{}
	decodeHooks := []viper.DecoderConfigOption{
		// this allows arrays to be decoded from env var separated by ",", example: MY_VAR="value1,value2,value3"
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
	}
	if err = v.Unmarshal(&dstConf, decodeHooks...); err != nil {
		log.Errorf("failed to unmarshal config: %v error: %v", value, err)
		return nil, err
	}
	return &dstConf, nil
}

const (
	// Halt is the key for l2gaspricer halt
	Halt         = "Halt"
	maxHaltDelay = 20
)

func (c *Client) fireHalt(key string, value *storage.ConfigChange) {
	switch key {
	case Halt:
		if value.OldValue.(string) != value.NewValue.(string) {
			random, _ := rand.Int(rand.Reader, big.NewInt(maxHaltDelay))
			delay := time.Second * time.Duration(random.Int64())
			log.Infof("halt changed from %s to %s delay halt %v", value.OldValue.(string), value.NewValue.(string), delay)
			time.Sleep(delay)
			os.Exit(1)
		}
	}
}
