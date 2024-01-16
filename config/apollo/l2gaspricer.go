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
	c.updateL2GasPricer(&c.config.L2GasPriceSuggester, newConf.L2GasPriceSuggester)
}

func (c *Client) updateL2GasPricer(dstConfig *gasprice.Config, srcConfig gasprice.Config) {
	if c == nil || !c.config.Apollo.Enable || dstConfig == nil {
		log.Infof("apollo is not enabled %v %v %v", c, dstConfig, srcConfig)
		return
	}
	if dstConfig.DefaultGasPriceWei != srcConfig.DefaultGasPriceWei {
		log.Infof("l2gaspricer default gas price changed from %d to %d",
			dstConfig.DefaultGasPriceWei, srcConfig.DefaultGasPriceWei)
		dstConfig.DefaultGasPriceWei = srcConfig.DefaultGasPriceWei
	}
	if dstConfig.MaxGasPriceWei != srcConfig.MaxGasPriceWei {
		log.Infof("l2gaspricer max gas price changed from %d to %d",
			dstConfig.MaxGasPriceWei, srcConfig.MaxGasPriceWei)
		dstConfig.MaxGasPriceWei = srcConfig.MaxGasPriceWei
	}
	if dstConfig.Factor != srcConfig.Factor {
		log.Infof("l2gaspricer factor changed from %v to %v",
			dstConfig.Factor, srcConfig.Factor)
		dstConfig.Factor = srcConfig.Factor
	}
	if dstConfig.GasPriceUsdt != srcConfig.GasPriceUsdt {
		log.Infof("l2gaspricer gas price usdt changed from %v to %v",
			dstConfig.GasPriceUsdt, srcConfig.GasPriceUsdt)
		dstConfig.GasPriceUsdt = srcConfig.GasPriceUsdt
	}
}

// FetchL2GasPricerConfig fetches the l2gaspricer config, called from gasprice module
func (c *Client) FetchL2GasPricerConfig(config *gasprice.Config) {
	if c == nil || !c.config.Apollo.Enable || config == nil {
		return
	}

	c.updateL2GasPricer(config, c.config.L2GasPriceSuggester)
}
