package apollo

import (
	"github.com/0xPolygonHermez/zkevm-node/gasprice"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/apolloconfig/agollo/v4/storage"
)

func (c *Client) loadL2GasPricer(value interface{}) {
	dstConf, err := c.unmarshal(value)
	if err != nil {
		log.Fatalf("failed to unmarshal l2gaspricer config: %v", err)
	}
	c.config.L2GasPriceSuggester = dstConf.L2GasPriceSuggester
	log.Infof("loaded l2gaspricer from apollo config: %+v", value.(string))
}

// fireL2GasPricer fires the l2gaspricer config change
// DefaultGasPriceWei
// MaxGasPriceWei
// Factor
// GasPriceUsdt
func (c *Client) fireL2GasPricer(key string, value *storage.ConfigChange) {
	newConf, err := c.unmarshal(value.NewValue)
	if err != nil {
		log.Errorf("failed to unmarshal l2gaspricer config: %v error: %v", value.NewValue, err)
		return
	}
	log.Infof("apollo l2gaspricer old config : %+v", value.OldValue.(string))
	log.Infof("apollo l2gaspricer config changed: %+v", value.NewValue.(string))
	gasprice.UpdateConfig(newConf.L2GasPriceSuggester)
}
