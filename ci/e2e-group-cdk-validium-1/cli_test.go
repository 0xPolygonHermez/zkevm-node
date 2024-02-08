package e2e

import (
	"context"
	"os/exec"
	"regexp"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygondatacommittee"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestSetDataAvailabilityProtocol(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	err := operations.Teardown()
	require.NoError(t, err)

	opsCfg := operations.GetDefaultOperationsConfig()

	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)

	err = opsman.Setup()
	require.NoError(t, err)

	clientL1, err := ethclient.Dial(operations.DefaultL1NetworkURL)
	require.NoError(t, err)

	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	require.NoError(t, err)

	zkEVM, err := polygonzkevm.NewPolygonzkevm(
		common.HexToAddress(operations.DefaultL1ZkEVMSmartContract),
		clientL1,
	)
	require.NoError(t, err)

	currentDAPAddr, err := zkEVM.DataAvailabilityProtocol(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	require.Equal(t, common.HexToAddress(operations.DefaultL1DataCommitteeContract), currentDAPAddr)

	// New DAC Setup
	newDAPAddr, tx, newDA, err := polygondatacommittee.DeployPolygondatacommittee(auth, clientL1)
	require.NoError(t, err)
	require.NotEqual(t, newDAPAddr, currentDAPAddr)
	require.NoError(t, operations.WaitTxToBeMined(ctx, clientL1, tx, operations.DefaultTimeoutTxToBeMined))

	tx, err = newDA.Initialize(auth)
	require.NoError(t, err)
	require.NoError(t, operations.WaitTxToBeMined(ctx, clientL1, tx, operations.DefaultTimeoutTxToBeMined))

	cmd := exec.Command("docker", "exec", "zkevm-sequence-sender",
		"/app/zkevm-node", "set-dap",
		"--da-addr", newDAPAddr.String(),
		"--network", "custom",
		"--custom-network-file", "/app/genesis.json",
		"--key-store-path", "/pk/sequencer.keystore",
		"--pw", "testonly",
		"--cfg", "/app/config.toml")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	txHash := common.HexToHash(extractHexFromString(string(output)))
	receipt, err := operations.WaitTxReceipt(ctx, txHash, operations.DefaultTimeoutTxToBeMined, clientL1)
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	currentDAPAddr, err = zkEVM.DataAvailabilityProtocol(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	require.Equal(t, newDAPAddr, currentDAPAddr)
}

func extractHexFromString(output string) string {
	re := regexp.MustCompile(`Transaction to set new data availability protocol sent. Hash: (0x[0-9a-fA-F]+)`)
	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
