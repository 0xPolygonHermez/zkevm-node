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

	bridgeDepositReceiverAddress    = "0xc949254d682d8c9ad5682521675b8f43b102aec4"
	bridgeDepositReceiverPrivateKey = "0xdfd01798f92667dbf91df722434e8fbe96af0211d4d1b82bbbbc8f1def7a814f"

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
	chkErr(err)

	// Eth client
	log.Infof("Connecting to l1")
	clientL1, err := ethclient.Dial(l1NetworkURL)
	chkErr(err)

	// Hermez client
	log.Infof("Connecting to l1")
	clientL2, err := ethclient.Dial(l2NetworkURL)
	chkErr(err)

	// Get network chain id
	log.Infof("Getting chainID L1")
	chainIDL1, err := clientL1.NetworkID(ctx)
	chkErr(err)

	// Preparing l1 acc info
	log.Infof("Creating deployer authorization")
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(l1AccHexPrivateKey, "0x"))
	chkErr(err)
	authDeployer, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDL1)
	chkErr(err)

	// Create sequencer auth
	log.Infof("Creating sequencer authorization")
	privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(sequencerPrivateKey, "0x"))
	chkErr(err)
	authSequencer, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDL1)
	chkErr(err)

	// Getting l1 info
	log.Infof("Getting L1 info")
	gasPrice, err := clientL1.SuggestGasPrice(ctx)
	chkErr(err)

	// Send some Ether from l1Acc to sequencer acc
	log.Infof("Transferring ETH to the sequencer")
	fromAddress := common.HexToAddress(l1AccHexAddress)
	nonce, err := clientL1.PendingNonceAt(ctx, fromAddress)
	chkErr(err)
	const gasLimit = 21000
	toAddress := common.HexToAddress(sequencerAddress)
	ethAmount, _ := big.NewInt(0).SetString("200000000000000000000", encoding.Base10)
	tx := types.NewTransaction(nonce, toAddress, ethAmount, gasLimit, gasPrice, nil)
	signedTx, err := authDeployer.Signer(authDeployer.From, tx)
	chkErr(err)
	err = clientL1.SendTransaction(ctx, signedTx)
	chkErr(err)
	_, err = scripts.WaitTxToBeMined(clientL1, signedTx.Hash(), txTimeout)
	chkErr(err)

	// Create matic maticTokenSC sc instance
	log.Infof("Loading Matic token SC instance")
	maticTokenSC, err := matic.NewMatic(cfg.NetworkConfig.MaticAddr, clientL1)
	chkErr(err)

	// Send matic to sequencer
	log.Infof("Transferring MATIC tokens to sequencer")
	maticAmount, _ := big.NewInt(0).SetString("200000000000000000000000", encoding.Base10)
	tx, err = maticTokenSC.Transfer(authDeployer, toAddress, maticAmount)
	chkErr(err)

	// wait matic transfer to be mined
	_, err = scripts.WaitTxToBeMined(clientL1, tx.Hash(), txTimeout)
	chkErr(err)

	// approve tokens to be used by PoE SC on behalf of the sequencer
	log.Infof("Approving tokens to be used by PoE on behalf of the sequencer")
	tx, err = maticTokenSC.Approve(authSequencer, cfg.NetworkConfig.PoEAddr, maticAmount)
	chkErr(err)
	_, err = scripts.WaitTxToBeMined(clientL1, tx.Hash(), txTimeout)
	chkErr(err)

	// Register the sequencer
	log.Infof("Registering the sequencer")
	ethermanConfig := etherman.Config{
		URL: l1NetworkURL,
	}
	etherman, err := etherman.NewClient(ethermanConfig, authSequencer, cfg.NetworkConfig.PoEAddr, cfg.NetworkConfig.MaticAddr)
	chkErr(err)
	tx, err = etherman.RegisterSequencer(l2NetworkURL)
	chkErr(err)

	// Wait sequencer to be registered
	log.Infof("waiting sequencer to be registered")
	_, err = scripts.WaitTxToBeMined(clientL1, tx.Hash(), txTimeout)
	chkErr(err)
	log.Infof("sequencer registered")

	// Deposit funds to L2 via bridge
	log.Infof("Depositing funds to L2 via bridge")
	balance, err := clientL1.BalanceAt(ctx, authSequencer.From, nil)
	chkErr(err)
	log.Debugf("ETH Balance of %v: %v", authSequencer.From.Hex(), balance.Text(encoding.Base10))

	const destNetwork = uint32(1)
	depositAmount, _ := big.NewInt(0).SetString("10000000000000000000", encoding.Base10)
	ethAddr := common.Address{}
	destAddr := common.HexToAddress(bridgeDepositReceiverAddress)
	sendL1Deposit(ctx, authDeployer, clientL1, ethAddr, depositAmount, destNetwork, &destAddr)

	lastBatchNumber, err := clientL2.BlockNumber(ctx)
	chkErr(err)

	// Proposing empty batch to trigger the l2 synchronization process
	forceBatchProposal(ctx, authSequencer, clientL1, cfg.NetworkConfig)

	expectedLastBatchNumber := lastBatchNumber + 1
	for {
		currentLastBatchNumber, err := clientL2.BlockNumber(ctx)
		chkErr(err)
		log.Infof("Waiting synchronizer to sync the forced empty batch. Current: %v Expected: %v", currentLastBatchNumber, expectedLastBatchNumber)
		if currentLastBatchNumber == expectedLastBatchNumber {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Claiming the funds deposited via bridge on L2
	deposit := deposit{
		TokenAddr:  common.Address{},
		Amount:     depositAmount,
		OrigNet:    0,
		DestNet:    destNetwork,
		DestAddr:   authSequencer.From,
		DepositCnt: 1,
	}
	smtProof := [][32]byte{
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0xad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5"),
		common.HexToHash("0xb4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d30"),
		common.HexToHash("0x21ddb9a356815c3fac1026b6dec5df3124afbadb485c9ba5a3e3398a04b7ba85"),
		common.HexToHash("0xe58769b32a1beaf1ea27375a44095a0d1fb664ce2dd358e7fcbfb78c26a19344"),
		common.HexToHash("0x0eb01ebfc9ed27500cd4dfc979272d1f0913cc9f66540d7e8005811109e1cf2d"),
		common.HexToHash("0x887c22bd8750d34016ac3c66b5ff102dacdd73f6b014e710b51e8022af9a1968"),
		common.HexToHash("0xffd70157e48063fc33c97a050f7f640233bf646cc98d9524c6b92bcf3ab56f83"),
		common.HexToHash("0x9867cc5f7f196b93bae1e27e6320742445d290f2263827498b54fec539f756af"),
		common.HexToHash("0xcefad4e508c098b9a7e1d8feb19955fb02ba9675585078710969d3440f5054e0"),
		common.HexToHash("0xf9dc3e7fe016e050eff260334f18a5d4fe391d82092319f5964f2e2eb7c1c3a5"),
		common.HexToHash("0xf8b13a49e282f609c317a833fb8d976d11517c571d1221a265d25af778ecf892"),
		common.HexToHash("0x3490c6ceeb450aecdc82e28293031d10c7d73bf85e57bf041a97360aa2c5d99c"),
		common.HexToHash("0xc1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb"),
		common.HexToHash("0x5c67add7c6caf302256adedf7ab114da0acfe870d449a3a489f781d659e8becc"),
		common.HexToHash("0xda7bce9f4e8618b6bd2f4132ce798cdc7a60e7e1460a7299e3c6342a579626d2"),
		common.HexToHash("0x2733e50f526ec2fa19a22b31e8ed50f23cd1fdf94c9154ed3a7609a2f1ff981f"),
		common.HexToHash("0xe1d3b5c807b281e4683cc6d6315cf95b9ade8641defcb32372f1c126e398ef7a"),
		common.HexToHash("0x5a2dce0a8a7f68bb74560f8f71837c2c2ebbcbf7fffb42ae1896f13f7c7479a0"),
		common.HexToHash("0xb46a28b6f55540f89444f63de0378e3d121be09e06cc9ded1c20e65876d36aa0"),
		common.HexToHash("0xc65e9645644786b620e2dd2ad648ddfcbf4a7e5b1a3a4ecfe7f64667a3f0b7e2"),
		common.HexToHash("0xf4418588ed35a2458cffeb39b93d26f18d2ab13bdce6aee58e7b99359ec2dfd9"),
		common.HexToHash("0x5a9c16dc00d6ef18b7933a6f8dc65ccb55667138776f7dea101070dc8796e377"),
		common.HexToHash("0x4df84f40ae0c8229d0d6069e5c8f39a7c299677a09d367fc7b05e3bc380ee652"),
		common.HexToHash("0xcdc72595f74c7b1043d0e1ffbab734648c838dfb0527d971b602bc216c9619ef"),
		common.HexToHash("0x0abf5ac974a1ed57f4050aa510dd9c74f508277b39d7973bb2dfccc5eeb0618d"),
		common.HexToHash("0xb8cd74046ff337f0a7bf2c8e03e10f642c1886798d71806ab1e888d9e5ee87d0"),
		common.HexToHash("0x838c5655cb21c6cb83313b5a631175dff4963772cce9108188b34ac87c81c41e"),
		common.HexToHash("0x662ee4dd2dd7b2bc707961b1e646c4047669dcb6584f0d8d770daf5d7e7deb2e"),
		common.HexToHash("0x388ab20e2573d171a88108e79d820e98f26c0b84aa8b2f4aa4968dbb818ea322"),
		common.HexToHash("0x93237c50ba75ee485f4c22adf2f741400bdf8d6a9cc7df7ecae576221665d735"),
		common.HexToHash("0x8448818bb4ae4562849e949e17ac16e0be16688e156b5cf15e098c627c0056a9"),
	}
	globalExitRoot := &globalExitRoot{
		GlobalExitRootNum: big.NewInt(1),
		ExitRoots: []common.Hash{
			common.HexToHash("0x843cb84814162b93794ad9087a037a1948f9aff051838ba3a93db0ac92b9f719"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
	}

	// Preparing bridge receiver acc info
	log.Infof("Creating bridge receiver authorization")
	privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(bridgeDepositReceiverPrivateKey, "0x"))
	chkErr(err)
	authBridgeReceiver, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDL1)
	chkErr(err)

	sendL2Claim(ctx, authBridgeReceiver, clientL2, deposit, smtProof, globalExitRoot)

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
	chkErr(err)

	tx, err := br.Bridge(auth, tokenAddr, amount, destNetwork, *destAddr)
	chkErr(err)

	log.Infof("Waiting L1Deposit to be mined")
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	chkErr(err)
	log.Infof("L1Deposit mined: %v", tx.Hash())
}

func forceBatchProposal(ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, cfg config.NetworkConfig) {
	log.Infof("Forcing batch proposal")

	poe, err := proofofefficiency.NewProofofefficiency(cfg.PoEAddr, client)
	chkErr(err)

	maticAmount, err := poe.CalculateSequencerCollateral(&bind.CallOpts{Pending: false})
	chkErr(err)
	log.Infof("Collateral: %v", maticAmount.Text(encoding.Base10))

	tx, err := poe.SendBatch(auth, []byte{}, maticAmount)
	chkErr(err)

	log.Infof("Waiting force batch proposal to be mined")
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	chkErr(err)
}

func sendL2Claim(ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, dep deposit, smtProof [][32]byte, globalExitRoot *globalExitRoot) {
	auth.GasPrice = big.NewInt(0)
	br, err := bridge.NewBridge(common.HexToAddress(l2BridgeAddr), client)
	chkErr(err)

	tx, err := br.Claim(auth, dep.TokenAddr, dep.Amount, dep.OrigNet, dep.DestNet,
		dep.DestAddr, smtProof, dep.DepositCnt, globalExitRoot.GlobalExitRootNum,
		globalExitRoot.ExitRoots[0], globalExitRoot.ExitRoots[1])
	chkErr(err)

	log.Infof("waiting L2 Claim tx to be mined")
	const txTimeout = 15 * time.Second
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	chkErr(err)

	log.Infof("wait for the consolidation")
	const timeToWaitForTheConsolidation = 30 * time.Second
	time.Sleep(timeToWaitForTheConsolidation)
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
