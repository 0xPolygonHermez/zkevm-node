package main

import (
	"context"
	"flag"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/bridge"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/matic"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/scripts"
	"github.com/hermeznetwork/hermez-core/test/operations"
	"github.com/urfave/cli/v2"
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

	txTimeout = 5 * time.Second
)

type deposit struct {
	TokenAddr  common.Address
	Amount     *big.Int
	OrigNet    uint32
	DestNet    uint32
	DestAddr   common.Address
	DepositCnt uint32
}

type globalExitRoot struct {
	BlockID           uint64
	BlockNumber       uint64
	GlobalExitRootNum *big.Int
	ExitRoots         []common.Hash
}

// TestStateTransition tests state transitions using the vector
func main() {
	ctx := context.Background()

	app := cli.NewApp()
	var n string
	flag.StringVar(&n, "network", "local", "")
	context := cli.NewContext(app, flag.CommandLine, nil)

	cfg, err := config.Load(context)
	checkErr(err)

	// Eth client
	log.Infof("Connecting to l1")
	client, err := ethclient.Dial(l1NetworkURL)
	checkErr(err)

	// Get network chain id
	log.Infof("Getting chainID")
	chainID, err := client.NetworkID(ctx)
	checkErr(err)

	// Preparing l1 acc info
	log.Infof("Preparing authorization")
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(l1AccHexPrivateKey, "0x"))
	checkErr(err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	checkErr(err)

	// Getting l1 info
	log.Infof("Getting L1 info")
	gasPrice, err := client.SuggestGasPrice(ctx)
	checkErr(err)

	// Send some Ether from l1Acc to sequencer acc
	log.Infof("Transferring ETH to the sequencer")
	fromAddress := common.HexToAddress(l1AccHexAddress)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	checkErr(err)
	const gasLimit = 21000
	toAddress := common.HexToAddress(sequencerAddress)
	ethAmount, _ := big.NewInt(0).SetString("200000000000000000000", encoding.Base10)
	tx := types.NewTransaction(nonce, toAddress, ethAmount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	checkErr(err)
	err = client.SendTransaction(ctx, signedTx)
	checkErr(err)
	_, err = scripts.WaitTxToBeMined(client, signedTx.Hash(), txTimeout)
	checkErr(err)

	// Create matic maticTokenSC sc instance
	log.Infof("Loading Matic token SC instance")
	maticTokenSC, err := operations.NewToken(cfg.NetworkConfig.MaticAddr, client)
	checkErr(err)

	// Send matic to sequencer
	log.Infof("Transferring MATIC tokens to sequencer")
	maticAmount, _ := big.NewInt(0).SetString("200000000000000000000000", encoding.Base10)
	tx, err = maticTokenSC.Transfer(auth, toAddress, maticAmount)
	checkErr(err)

	// wait matic transfer to be mined
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	checkErr(err)

	// Create sequencer auth
	log.Infof("Creating sequencer authorization")
	privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(sequencerPrivateKey, "0x"))
	checkErr(err)
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	checkErr(err)

	// approve tokens to be used by PoE SC on behalf of the sequencer
	log.Infof("Approving tokens to be used by PoE on behalf of the sequencer")
	tx, err = maticTokenSC.Approve(auth, cfg.NetworkConfig.PoEAddr, maticAmount)
	checkErr(err)
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	checkErr(err)

	// Register the sequencer
	log.Infof("Registering the sequencer")
	ethermanConfig := etherman.Config{
		URL: l1NetworkURL,
	}
	etherman, err := etherman.NewClient(ethermanConfig, auth, cfg.NetworkConfig.PoEAddr, cfg.NetworkConfig.MaticAddr)
	checkErr(err)
	tx, err = etherman.RegisterSequencer(l2NetworkURL)
	checkErr(err)

	// Wait sequencer to be registered
	log.Infof("waiting sequencer to be registered")
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	checkErr(err)
	log.Infof("sequencer registered")

	// Deposit funds to L2 via bridge
	log.Infof("Depositing funds to L2 via bridge")
	balance, err := client.BalanceAt(ctx, auth.From, nil)
	checkErr(err)
	log.Debugf("ETH Balance of %v: %v", auth.From.Hex(), balance.Text(encoding.Base10))

	const destNetwork = uint32(1)
	depositAmount, _ := big.NewInt(0).SetString("10000000000000000000", encoding.Base10)
	ethAddr := common.Address{}
	destAddr := common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4")
	// pk 0xdfd01798f92667dbf91df722434e8fbe96af0211d4d1b82bbbbc8f1def7a814f
	sendL1Deposit(ctx, auth, client, ethAddr, depositAmount, destNetwork, &destAddr)

	// Proposing empty batch to trigger the l2 synchronization process
	forceBatchProposal(ctx, auth, client, cfg.NetworkConfig)

	// Claiming the funds deposited via bridge on L2
	// TODO: Get the values for deposit, smtProof and globalExitRoot
	deposit := deposit{
		TokenAddr:  common.Address{},
		Amount:     depositAmount,
		OrigNet:    0,
		DestNet:    destNetwork,
		DestAddr:   auth.From,
		DepositCnt: 1,
	}
	smtProof := [][32]byte{
		hexAddressTo32bytes("0x0000000000000000000000000000000000000000000000000000000000000000"),
		hexAddressTo32bytes("0xad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5"),
		hexAddressTo32bytes("0xb4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d30"),
		hexAddressTo32bytes("0x21ddb9a356815c3fac1026b6dec5df3124afbadb485c9ba5a3e3398a04b7ba85"),
		hexAddressTo32bytes("0xe58769b32a1beaf1ea27375a44095a0d1fb664ce2dd358e7fcbfb78c26a19344"),
		hexAddressTo32bytes("0x0eb01ebfc9ed27500cd4dfc979272d1f0913cc9f66540d7e8005811109e1cf2d"),
		hexAddressTo32bytes("0x887c22bd8750d34016ac3c66b5ff102dacdd73f6b014e710b51e8022af9a1968"),
		hexAddressTo32bytes("0xffd70157e48063fc33c97a050f7f640233bf646cc98d9524c6b92bcf3ab56f83"),
		hexAddressTo32bytes("0x9867cc5f7f196b93bae1e27e6320742445d290f2263827498b54fec539f756af"),
		hexAddressTo32bytes("0xcefad4e508c098b9a7e1d8feb19955fb02ba9675585078710969d3440f5054e0"),
		hexAddressTo32bytes("0xf9dc3e7fe016e050eff260334f18a5d4fe391d82092319f5964f2e2eb7c1c3a5"),
		hexAddressTo32bytes("0xf8b13a49e282f609c317a833fb8d976d11517c571d1221a265d25af778ecf892"),
		hexAddressTo32bytes("0x3490c6ceeb450aecdc82e28293031d10c7d73bf85e57bf041a97360aa2c5d99c"),
		hexAddressTo32bytes("0xc1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb"),
		hexAddressTo32bytes("0x5c67add7c6caf302256adedf7ab114da0acfe870d449a3a489f781d659e8becc"),
		hexAddressTo32bytes("0xda7bce9f4e8618b6bd2f4132ce798cdc7a60e7e1460a7299e3c6342a579626d2"),
		hexAddressTo32bytes("0x2733e50f526ec2fa19a22b31e8ed50f23cd1fdf94c9154ed3a7609a2f1ff981f"),
		hexAddressTo32bytes("0xe1d3b5c807b281e4683cc6d6315cf95b9ade8641defcb32372f1c126e398ef7a"),
		hexAddressTo32bytes("0x5a2dce0a8a7f68bb74560f8f71837c2c2ebbcbf7fffb42ae1896f13f7c7479a0"),
		hexAddressTo32bytes("0xb46a28b6f55540f89444f63de0378e3d121be09e06cc9ded1c20e65876d36aa0"),
		hexAddressTo32bytes("0xc65e9645644786b620e2dd2ad648ddfcbf4a7e5b1a3a4ecfe7f64667a3f0b7e2"),
		hexAddressTo32bytes("0xf4418588ed35a2458cffeb39b93d26f18d2ab13bdce6aee58e7b99359ec2dfd9"),
		hexAddressTo32bytes("0x5a9c16dc00d6ef18b7933a6f8dc65ccb55667138776f7dea101070dc8796e377"),
		hexAddressTo32bytes("0x4df84f40ae0c8229d0d6069e5c8f39a7c299677a09d367fc7b05e3bc380ee652"),
		hexAddressTo32bytes("0xcdc72595f74c7b1043d0e1ffbab734648c838dfb0527d971b602bc216c9619ef"),
		hexAddressTo32bytes("0x0abf5ac974a1ed57f4050aa510dd9c74f508277b39d7973bb2dfccc5eeb0618d"),
		hexAddressTo32bytes("0xb8cd74046ff337f0a7bf2c8e03e10f642c1886798d71806ab1e888d9e5ee87d0"),
		hexAddressTo32bytes("0x838c5655cb21c6cb83313b5a631175dff4963772cce9108188b34ac87c81c41e"),
		hexAddressTo32bytes("0x662ee4dd2dd7b2bc707961b1e646c4047669dcb6584f0d8d770daf5d7e7deb2e"),
		hexAddressTo32bytes("0x388ab20e2573d171a88108e79d820e98f26c0b84aa8b2f4aa4968dbb818ea322"),
		hexAddressTo32bytes("0x93237c50ba75ee485f4c22adf2f741400bdf8d6a9cc7df7ecae576221665d735"),
		hexAddressTo32bytes("0x8448818bb4ae4562849e949e17ac16e0be16688e156b5cf15e098c627c0056a9"),
	}
	globalExitRoot := &globalExitRoot{
		GlobalExitRootNum: big.NewInt(1),
		ExitRoots: []common.Hash{
			common.HexToHash("0x843cb84814162b93794ad9087a037a1948f9aff051838ba3a93db0ac92b9f719"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
	}
	sendL2Claim(ctx, auth, client, deposit, smtProof, globalExitRoot)

	log.Infof("Network initialized properly")
}

// sendL1Deposit sends a deposit from l1 to l2
func sendL1Deposit(ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, tokenAddr common.Address, amount *big.Int,
	destNetwork uint32, destAddr *common.Address,
) {
	emptyAddr := common.Address{}
	if tokenAddr == emptyAddr {
		auth.Value = amount
	}
	if destAddr == nil {
		destAddr = &auth.From
	}
	br, err := bridge.NewBridge(common.HexToAddress(l1BridgeAddr), client)
	checkErr(err)

	tx, err := br.Bridge(auth, tokenAddr, amount, destNetwork, *destAddr)
	checkErr(err)

	log.Infof("Waiting L1Deposit to be mined")
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	checkErr(err)
	log.Infof("L1Deposit mined: %v", tx.Hash())
}

func forceBatchProposal(ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, cfg config.NetworkConfig) {
	log.Infof("Forcing batch proposal")

	poe, err := proofofefficiency.NewProofofefficiency(cfg.PoEAddr, client)
	checkErr(err)

	maticAmount, err := poe.CalculateSequencerCollateral(&bind.CallOpts{Pending: false})
	checkErr(err)
	log.Infof("Collateral: %v", maticAmount.Text(encoding.Base10))

	m, err := matic.NewMatic(cfg.MaticAddr, client)
	checkErr(err)
	balance, err := m.BalanceOf(nil, auth.From)
	checkErr(err)

	log.Infof("MATIC Balance of %v: %v", auth.From.Hex(), balance.Text(encoding.Base10))

	tx, err := poe.SendBatch(auth, []byte{}, maticAmount)
	checkErr(err)

	log.Infof("Waiting force batch proposal to be mined")
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	checkErr(err)
}

func sendL2Claim(ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, dep deposit, smtProof [][32]byte, globalExitRoot *globalExitRoot) {
	auth.GasPrice = big.NewInt(0)
	br, err := bridge.NewBridge(common.HexToAddress(l2BridgeAddr), client)
	checkErr(err)

	tx, err := br.Claim(auth, dep.TokenAddr, dep.Amount, dep.OrigNet, dep.DestNet,
		dep.DestAddr, smtProof, dep.DepositCnt, globalExitRoot.GlobalExitRootNum,
		globalExitRoot.ExitRoots[0], globalExitRoot.ExitRoots[1])
	checkErr(err)

	log.Infof("waiting L2 Claim tx to be mined")
	const txTimeout = 15 * time.Second
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	checkErr(err)

	log.Infof("wait for the consolidation")
	const timeToWaitForTheConsolidation = 30 * time.Second
	time.Sleep(timeToWaitForTheConsolidation)
}

func hexAddressTo32bytes(hex string) [32]byte {
	addr := common.HexToAddress(hex)
	addrBytes := addr.Bytes()
	result := [32]byte{}
	copy(result[:], addrBytes[:32])
	return result
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
