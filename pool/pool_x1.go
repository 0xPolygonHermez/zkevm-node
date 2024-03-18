package pool

import (
	"context"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// BridgeClaimMethodSignature for tracking BridgeClaimMethodSignature method
	BridgeClaimMethodSignature = "0xccaa2d11"
)

func contains(s []string, ele common.Address) bool {
	for _, e := range s {
		if common.HexToAddress(e) == ele {
			return true
		}
	}
	return false
}

// StartRefreshingWhiteAddressesPeriodically will make this instance of the pool
// to check periodically(accordingly to the configuration) for updates regarding
// the white address and update the in memory blocked addresses
func (p *Pool) StartRefreshingWhiteAddressesPeriodically() {
	p.refreshWhitelistedAddresses()
	go func(p *Pool) {
		for {
			time.Sleep(p.cfg.IntervalToRefreshWhiteAddresses.Duration)
			p.refreshWhitelistedAddresses()
		}
	}(p)
}

// refreshWhitelistedAddresses refreshes the list of whitelisted addresses for the provided instance of pool
func (p *Pool) refreshWhitelistedAddresses() {
	whitelistedAddresses, err := p.storage.GetAllAddressesWhitelisted(context.Background())
	if err != nil {
		log.Error("failed to load whitelisted addresses")
		return
	}

	whitelistedAddressesMap := sync.Map{}
	for _, whitelistedAddress := range whitelistedAddresses {
		whitelistedAddressesMap.Store(whitelistedAddress.String(), 1)
		p.whitelistedAddresses.Store(whitelistedAddress.String(), 1)
	}

	nonWhitelistedAddresses := []string{}
	p.whitelistedAddresses.Range(func(key, value any) bool {
		addrHex := key.(string)
		_, found := whitelistedAddressesMap.Load(addrHex)
		if found {
			return true
		}

		nonWhitelistedAddresses = append(nonWhitelistedAddresses, addrHex)
		return true
	})

	for _, nonWhitelistedAddress := range nonWhitelistedAddresses {
		p.whitelistedAddresses.Delete(nonWhitelistedAddress)
	}
}
