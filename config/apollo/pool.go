package apollo

import (
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/apolloconfig/agollo/v4/storage"
)

// loadPool loads the pool config from apollo
func (c *Client) loadPool(value interface{}) {
	dstConf, err := c.unmarshal(value)
	if err != nil {
		log.Fatalf("failed to unmarshal pool config: %v", err)
	}
	c.config.Pool = dstConf.Pool
	log.Infof("loaded pool from apollo config: %+v", value.(string))
}

// firePool fires the pool config change
// AccountQueue
// GlobalQueue
// FreeGasAddress
func (c *Client) firePool(key string, value *storage.ConfigChange) {
	newConf, err := c.unmarshal(value.NewValue)
	if err != nil {
		log.Errorf("failed to unmarshal pool config: %v error: %v", value.NewValue, err)
		return
	}
	log.Infof("apollo pool old config : %+v", value.OldValue.(string))
	log.Infof("apollo pool config changed: %+v", value.NewValue.(string))

	pool.UpdateConfig(newConf.Pool)
}
