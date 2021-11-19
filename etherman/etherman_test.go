package etherman

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeOneTxData(t *testing.T) {
	dHex := "f8638084773594008208349490f79bf6eb2c4f870365e785982e1f101e93b906808026a08c3c762704595093e433b4cf9c36be5e5d2ac61b8e05e6e90011a3113544dee6a00e08c86807bc043893663af53c9a33394aee2cb427c0add182a38e57eb2c0411"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	tx, err := decodeTxs(data, big.NewInt(1))
	require.NoError(t, err)
	// fmt.Println(tx)
	var addr common.Address
	err = addr.UnmarshalText([]byte("0x90F79bf6EB2c4f870365E785982E1f101E93b906"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To)

	//This test should fail verifying the signature
	fmt.Println("###########################")
	dHex = "f86c088504a817c8008252089411111111111111111111111111111111111111118802c68af0bb1400008026a0758b52ed7380ef07d97a26904f6f2340e9437d3f44d4a950db48de846d18d6e5a0562ead2f0619ae253d65196b22026728829dd785c6a489cd2a546c6066d32c2a"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(1))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x1212121212121212121212121212121212121212"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To)
	fmt.Println("###########################")
}

func TestDecodeMultipleTxData(t *testing.T) {
	dHex := "f86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83f86c088504a817c8008252089411111111111111111111111111111111111111118802c68af0bb1400008026a08975cf0fe106a0396649d37c5292274f66193346ce5b07bcbaa8dc248f7f5496a0684963f42a662640b6d27e3552151e3f42744c9b4d72dd70b262ca94b9473c94f869028504a817c8008252089412121212121212121212121212121212121212128506fc23ac008025a0528b1dd150ccae6e83fcc44bff11928ca635f0fc6819836a14d526af1ecf0519a02a96710022671e44c81a6f19b88f605567d32dd97508aa84c830ec9d4a4aa0d2"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	txs, err := decodeTxs(data, big.NewInt(1))
	require.NoError(t, err)
	res := []string{"0x3535353535353535353535353535353535353535", "0x1111111111111111111111111111111111111111", "0x1212121212121212121212121212121212121212"}
	for k,tx := range txs {
		var addr common.Address
		err = addr.UnmarshalText([]byte(res[k]))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To)
	}
}