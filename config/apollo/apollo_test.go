package apollo

import (
	"testing"

	nodeConfig "github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

func TestApolloClient_LoadConfig(t *testing.T) {
	nc := &nodeConfig.Config{
		Apollo: types.ApolloConfig{
			IP:            "",
			AppID:         "x1-devnet",
			NamespaceName: "l2gaspricer.txt,l2gaspricerHalt.properties",
			Enable:        false,
		},
	}
	client := NewClient(nc)

	client.LoadConfig()
	t.Log(nc.L2GasPriceSuggester)
	// time.Sleep(2 * time.Minute)
	t.Log(nc.L2GasPriceSuggester)
}
