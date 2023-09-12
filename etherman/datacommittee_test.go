package etherman

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateDataCommitteeEvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, auth, _, _, da := newTestingEnv()

	// Update the committee
	requiredAmountOfSignatures := big.NewInt(2)
	URLs := []string{"1", "2", "3"}
	addrs := []common.Address{
		common.HexToAddress("0x1"),
		common.HexToAddress("0x2"),
		common.HexToAddress("0x3"),
	}
	addrsBytes := []byte{}
	for _, addr := range addrs {
		addrsBytes = append(addrsBytes, addr.Bytes()...)
	}
	_, err := da.SetupCommittee(auth, requiredAmountOfSignatures, URLs, addrsBytes)
	require.NoError(t, err)
	ethBackend.Commit()

	// Assert the committee update
	actualSetup, err := etherman.GetCurrentDataCommittee()
	require.NoError(t, err)
	expectedMembers := []DataCommitteeMember{}
	expectedSetup := DataCommittee{
		RequiredSignatures: uint64(len(URLs) - 1),
		AddressesHash:      crypto.Keccak256Hash(addrsBytes),
	}
	for i, url := range URLs {
		expectedMembers = append(expectedMembers, DataCommitteeMember{
			URL:  url,
			Addr: addrs[i],
		})
	}
	expectedSetup.Members = expectedMembers
	assert.Equal(t, expectedSetup, *actualSetup)
}
