package bridgetree

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckEmptyMerkleTree(t *testing.T) {
	expectedRoot := "27ae5ba08d7291c96c8cbddcc148bf48a6d68c7974b94356f53754ef6171d757"

	bt := NewBridgeTree(32)
	root := bt.GetRoot()
	assert.Equal(t, expectedRoot, hex.EncodeToString(root[:]))
}

func TestCheckMerkleTree(t *testing.T) {
	expectedRoot := "ef6323ce57f75effe515924a57e60e54c93d48a99a0492de40b6e911ca639895"

	bt := NewBridgeTree(32)
	leafValue, _ := formatBytes32String("1")
	bt.Add(leafValue)
	root := bt.GetRoot()
	assert.Equal(t, expectedRoot, hex.EncodeToString(root[:]))
	proof := bt.GetProofTreeByIndex(0)
	index := 0
	verification := VerifyMerkleProof(leafValue, proof, index, root)
	assert.True(t, verification)
}

func TestCheckAdd1LeafToTheMerkleTree(t *testing.T) {
	expectedRoot := "413f0936195da8c24ee159c80f1b8363ca6aa13ab22288e72a9212c6ae213a61"

	bt := NewBridgeTree(32)
	leafValue, _ := formatBytes32String("123")
	bt.Add(leafValue)
	root := bt.GetRoot()

	// verify root
	currentNode := leafValue
	for i := 0; i < int(bt.height); i++ {
		currentNode = hash(currentNode, bt.zeroHashes[i])
	}
	assert.Equal(t, currentNode, root)

	assert.Equal(t, expectedRoot, hex.EncodeToString(root[:]))

	// check merkle proof
	proof := bt.GetProofTreeByIndex(0)
	index := 0
	verification := VerifyMerkleProof(leafValue, proof, index, root)
	assert.True(t, verification)
}

func TestCheckAdd2LeavesToTheMerkleTree(t *testing.T) {
	expectedRoot := "fcd41bd90af952325c71c02b37baf7ab41fab57e1643893913ac69fa5268a64a"

	bt := NewBridgeTree(32)
	leafValue, _ := formatBytes32String("123")
	leafValue2, _ := formatBytes32String("456")

	bt.Add(leafValue)
	bt.Add(leafValue2)

	root := bt.GetRoot()

	// verify root
	currentNode := hash(leafValue, leafValue2)
	for i := 1; i < int(bt.height); i++ {
		currentNode = hash(currentNode, bt.zeroHashes[i])
	}
	assert.Equal(t, currentNode, root)

	assert.Equal(t, expectedRoot, hex.EncodeToString(root[:]))

	// check merkle proof
	index := 0
	proof := bt.GetProofTreeByIndex(index)
	verification := VerifyMerkleProof(leafValue, proof, index, root)
	assert.True(t, verification)

	// check merkle proof
	index2 := 1
	proof2 := bt.GetProofTreeByIndex(index2)
	verification2 := VerifyMerkleProof(leafValue2, proof2, index2, root)
	assert.True(t, verification2)

	assert.False(t, VerifyMerkleProof(leafValue, proof2, index2, root))
	assert.False(t, VerifyMerkleProof(leafValue, proof, index2, root))
	assert.False(t, VerifyMerkleProof(leafValue2, proof2, index, root))
	assert.False(t, VerifyMerkleProof(leafValue2, proof, index2, root))
	assert.False(t, VerifyMerkleProof(leafValue2, proof2, index2+1, root))
}

func TestCheckAddMoreLeavesToTheMerkleTree(t *testing.T) {
	expectedRoot1 := "eb44e3204a8800f57dc23e646f3bccbf57ea22344d989f8958aad4bceb2f4039"
	expectedRoot2 := "4bc24305b376a6724be01609a49af47ead97c3849ae728771f1ac578335c6bf9"

	bt := NewBridgeTree(32)
	leafValue, _ := formatBytes32String("123")
	leafValue2, _ := formatBytes32String("456")
	leafValue3, _ := formatBytes32String("1")
	leafValue4, _ := formatBytes32String("2")
	leafValue5, _ := formatBytes32String("3")
	leafValue6, _ := formatBytes32String("4")
	leafValue7, _ := formatBytes32String("5")

	bt.Add(leafValue)
	bt.Add(leafValue2)
	bt.Add(leafValue3)
	bt.Add(leafValue4)
	bt.Add(leafValue5)
	bt.Add(leafValue6)
	root := bt.GetRoot()

	// verify root
	assert.Equal(t, expectedRoot1, hex.EncodeToString(root[:]))

	bt.Add(leafValue7)
	root = bt.GetRoot()

	// verify root
	assert.Equal(t, expectedRoot2, hex.EncodeToString(root[:]))

	// check merkle proof
	index := 0
	proof := bt.GetProofTreeByIndex(index)
	verification := VerifyMerkleProof(leafValue, proof, index, root)
	assert.True(t, verification)

	// check merkle proof
	index2 := 1
	proof2 := bt.GetProofTreeByIndex(index2)
	verification2 := VerifyMerkleProof(leafValue2, proof2, index2, root)
	assert.True(t, verification2)
}

func formatBytes32String(text string) ([32]byte, error) {
	bText := []byte(text)
	if len(bText) > 31 {
		return [32]byte{}, fmt.Errorf("text is more than 31 bytes long")
	}
	var res [32]byte
	copy(res[:], bText)
	return res, nil
}

func TestFormatBytes32String(t *testing.T) {
	text := []string{"Hello World!", "1", "123", "456"}
	expectedHex := []string{
		"48656c6c6f20576f726c64210000000000000000000000000000000000000000",
		"3100000000000000000000000000000000000000000000000000000000000000",
		"3132330000000000000000000000000000000000000000000000000000000000",
		"3435360000000000000000000000000000000000000000000000000000000000",
	}
	for i := 0; i < len(text); i++ {
		expected, _ := hex.DecodeString(expectedHex[i])
		res, err := formatBytes32String(text[i])
		require.NoError(t, err)
		require.Equal(t, res[:], expected)
	}
}
