package etherman

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ptypes "github.com/0xPolygon/polygon-sdk/types"
)

func TestDecodeOneTxData(t *testing.T) {
	dHex := "f86c088504a817c8008252089411111111111111111111111111111111111111118802c68af0bb1400008026a0758b52ed7380ef07d97a26904f6f2340e9437d3f44d4a950db48de846d18d6e5a0562ead2f0619ae253d65196b22026728829dd785c6a489cd2a546c6066d32c2a"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	tx, err := decodeTxs(data, big.NewInt(1))
	require.NoError(t, err)
	// fmt.Println(tx)
	var addr ptypes.Address
	err = addr.UnmarshalText([]byte("0x1111111111111111111111111111111111111111"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To)
	err = addr.UnmarshalText([]byte("0xD3E442496EB66a4748912ec4A3b7A111d0B855d6"))
	require.NoError(t, err)
	assert.Equal(t, addr, tx[0].From)

	dHex = "f86b028504a817c8008252089412121212121212121212121212121212121212128506fc23ac0080820226a0a2402d3351e8ec9b0a221d7ff48aca682c646528b1558e74a0b558a943c0f2a3a05c7fff6db65e560833d7f85f07ed092fb63aaa2307b60cdf03b63a2effc6a340"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(257))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x1212121212121212121212121212121212121212"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To)
	err = addr.UnmarshalText([]byte("0x178513579470dc158FA85c2AeA86e18a8c955402"))
	require.NoError(t, err)
	assert.Equal(t, addr, tx[0].From)

	dHex = "f86e5a8502540be4008259d8941234123412341234123412341234123412341234880214e8348c4f000082123442a0c6ff1e0034458c8dbf64966f49031e44c6509f85545b49d4df2a953e9f4d1324a07403e62dda1922fb1e226632e21e7382c345377ff46e6a43b79f169570e5a725"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(15))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x1234123412341234123412341234123412341234"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To)
	err = addr.UnmarshalText([]byte("0xfF06ad5d076fa274B49C297f3fE9e29B5bA9AaDC"))
	require.NoError(t, err)
	assert.Equal(t, addr, tx[0].From)

	dHex = "f8701c8502540be400823a9894987698769876987698769876987698769876987688011c37937e0800008256788202e0a0ae3c16aaf6a780e085f5f919b0d1e5f07a1c014ed4d700c3ab189b4f98677f38a0316608701f846807f96f496b9994e7a821fc9322383e54f84bd235b93fb774b0"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(350))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x9876987698769876987698769876987698769876"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To)
	err = addr.UnmarshalText([]byte("0x96E2793020bbbCDD0256B08244e1b3dF8B6EEB8c"))
	require.NoError(t, err)
	assert.Equal(t, addr, tx[0].From)

	dHex = "f873528504a817c800824a389480808080808080808080808080808080808080808802c68af0bb140000851234567890820b13a06790f151dab1b65b577532b479938f45abf871d6978e98ee0db53331ad548709a021b5038ec4cf5f7edc3f2ebac9ec5fdaa602f48b53986111a9bddaf5ad5b771c"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(1400))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x8080808080808080808080808080808080808080"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To)
	err = addr.UnmarshalText([]byte("0x2e988A386a799F506693793c6A5AF6B54dfAaBfB"))
	require.NoError(t, err)
	assert.Equal(t, addr, tx[0].From)

	dHex = "f8702f85055ae826008261a8941111111111222222222233333333334444444444880cc47f20295c0000841122334427a01a3cf5ea05180ac59514dc64d91a1d615857cad466dacfbde4dd06f1988a1074a05ed158fd5c7d54bd827e94683c7f74dbd64be9bb69cbb3e092d317cc4758a146"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(2))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x1111111111222222222233333333334444444444"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To)
	err = addr.UnmarshalText([]byte("0x6A5575E1230543D19e368B13446a7082Be6E0B47"))
	require.NoError(t, err)
	assert.Equal(t, addr, tx[0].From)

}

func TestDecodeMultipleTxData(t *testing.T) {
	dHex := "f86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83f86c088504a817c8008252089411111111111111111111111111111111111111118802c68af0bb1400008026a08975cf0fe106a0396649d37c5292274f66193346ce5b07bcbaa8dc248f7f5496a0684963f42a662640b6d27e3552151e3f42744c9b4d72dd70b262ca94b9473c94f869028504a817c8008252089412121212121212121212121212121212121212128506fc23ac008025a0528b1dd150ccae6e83fcc44bff11928ca635f0fc6819836a14d526af1ecf0519a02a96710022671e44c81a6f19b88f605567d32dd97508aa84c830ec9d4a4aa0d2"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	txs, err := decodeTxs(data, big.NewInt(1))
	require.NoError(t, err)
	res := []string{"0x3535353535353535353535353535353535353535", "0x1111111111111111111111111111111111111111", "0x1212121212121212121212121212121212121212"}
	fromRes := []string{"0x9d8A62f656a8d1615C1294fd71e9CFb3E4855A4F","0x674b6Bb9B00Dd754094D3D6a129696D194F45FE7","0x674b6Bb9B00Dd754094D3D6a129696D194F45FE7"}
	for k,tx := range txs {
		var addr ptypes.Address
		err = addr.UnmarshalText([]byte(res[k]))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To)
		err = addr.UnmarshalText([]byte(fromRes[k]))
		require.NoError(t, err)
		assert.Equal(t, addr, tx.From)
	}
}