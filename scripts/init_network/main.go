package main

import (
	"context"
	"time"

	NW "github.com/hermeznetwork/hermez-core/tools/network"
)

const (
	l1NetworkURL = "http://localhost:8545"
	l2NetworkURL = "http://localhost:8123"

	l1BridgeAddr = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	l2BridgeAddr = "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988"

	l1AccHexAddress    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	l1AccHexPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	sequencerAddress    = "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"
	sequencerPrivateKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"

	bridgeDepositReceiverAddress    = "0xc949254d682d8c9ad5682521675b8f43b102aec4"
	bridgeDepositReceiverPrivateKey = "0xdfd01798f92667dbf91df722434e8fbe96af0211d4d1b82bbbbc8f1def7a814f"

	txTimeout = 60 * time.Second
)

func main() {
	ctx := context.Background()
	NW.InitNetwork(ctx,
		l1NetworkURL,
		l2NetworkURL,
		l1BridgeAddr,
		l2BridgeAddr,
		l1AccHexAddress,
		l1AccHexPrivateKey,
		sequencerAddress,
		sequencerPrivateKey,
		bridgeDepositReceiverAddress,
		bridgeDepositReceiverPrivateKey,
		txTimeout)

}
