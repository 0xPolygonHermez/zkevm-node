package main

// import (
// 	"context"
// 	"log"
// 	"time"

// 	NW "github.com/0xPolygonHermez/zkevm-node/tools/network"
// )

// func main() {
// 	ctx := context.Background()
// 	if err := NW.InitNetwork(ctx,
// 		NW.InitNetworkConfig{
// 			L1NetworkURL: "http://localhost:8545",
// 			L2NetworkURL: "http://localhost:8123",
// 			L1BridgeAddr: "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9",
// 			L2BridgeAddr: "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988",
// 			L1Deployer: NW.L1Deployer{
// 				Address:                  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
// 				PrivateKey:               "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
// 				L1ETHAmountToSequencer:   "200000000000000000000",
// 				L1PolAmountToSequencer: "200000000000000000000000",
// 			},
// 			sequencerAddress:    "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
// 			SequencerPrivateKey: "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e",
// 			TxTimeout:           time.Minute,
// 		}); err != nil {
// 		log.Fatal(err)
// 	}
// }
