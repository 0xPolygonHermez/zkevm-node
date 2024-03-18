package pool

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// apolloConfig is the apollo pool dynamic config
type apolloConfig struct {
	EnableApollo     bool
	FreeGasAddresses []string
	GlobalQueue      uint64
	AccountQueue     uint64
	EnableWhitelist  bool

	sync.RWMutex
}

var apolloConf = &apolloConfig{}

// getApolloConfig returns the pool singleton instance
func getApolloConfig() *apolloConfig {
	return apolloConf
}

// Enable returns true if apollo is enabled
func (c *apolloConfig) enable() bool {
	if c == nil || !c.EnableApollo {
		return false
	}
	c.RLock()
	defer c.RUnlock()
	return c.EnableApollo
}

func (c *apolloConfig) setFreeGasAddresses(freeGasAddrs []string) {
	if c == nil || !c.EnableApollo {
		return
	}
	c.FreeGasAddresses = make([]string, len(freeGasAddrs))
	copy(c.FreeGasAddresses, freeGasAddrs)
}

// UpdateConfig updates the apollo config
// GlobalQueue
// AccountQueue
// FreeGasAddress
// EnableWhitelist
func UpdateConfig(apolloConfig Config) {
	getApolloConfig().Lock()
	getApolloConfig().EnableApollo = true
	getApolloConfig().GlobalQueue = apolloConfig.GlobalQueue
	getApolloConfig().AccountQueue = apolloConfig.AccountQueue
	getApolloConfig().setFreeGasAddresses(apolloConfig.FreeGasAddress)
	getApolloConfig().EnableWhitelist = apolloConfig.EnableWhitelist
	getApolloConfig().Unlock()
}

func isFreeGasAddress(localFreeGasAddrs []string, address common.Address) bool {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return contains(getApolloConfig().FreeGasAddresses, address)
	}

	return contains(localFreeGasAddrs, address)
}

func getGlobalQueue(globalQueue uint64) uint64 {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().GlobalQueue
	}

	return globalQueue
}

func getAccountQueue(accountQueue uint64) uint64 {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().AccountQueue
	}

	return accountQueue
}

func getEnableWhitelist(enableWhitelist bool) bool {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().EnableWhitelist
	}

	return enableWhitelist
}
