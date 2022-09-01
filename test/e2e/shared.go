package e2e

import (
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/core/types"
)

var networks = []struct {
	Name       string
	URL        string
	ChainID    uint64
	PrivateKey string
}{
	{Name: "Local L1", URL: operations.DefaultL1NetworkURL, ChainID: operations.DefaultL1ChainID, PrivateKey: operations.DefaultSequencerPrivateKey},
	{Name: "Local L2", URL: operations.DefaultL2NetworkURL, ChainID: operations.DefaultL2ChainID, PrivateKey: operations.DefaultSequencerPrivateKey},
}

func logTx(tx *types.Transaction) {
	sender, _ := state.GetSender(*tx)
	log.Debugf("********************")
	log.Debugf("Hash: %v", tx.Hash())
	log.Debugf("From: %v", sender)
	log.Debugf("Nonce: %v", tx.Nonce())
	log.Debugf("ChainId: %v", tx.ChainId())
	log.Debugf("To: %v", tx.To())
	log.Debugf("Gas: %v", tx.Gas())
	log.Debugf("GasPrice: %v", tx.GasPrice())
	log.Debugf("Cost: %v", tx.Cost())

	// b, _ := tx.MarshalBinary()
	//log.Debugf("RLP: ", hex.EncodeToHex(b))
	log.Debugf("********************")
}
