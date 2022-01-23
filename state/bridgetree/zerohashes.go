package bridgetree

// HashZero is an empty hash
var HashZero = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// ZeroHashes is an array of calculated zero hashes for each level of the tree
var ZeroHashes [][32]byte

func init() {
	ZeroHashes = generateZeroHashes(32)
}

func generateZeroHashes(height uint8) [][32]byte {
	var zeroHashes = [][32]byte{
		HashZero,
	}
	for i := 1; i <= int(height); i++ {
		zeroHashes = append(zeroHashes, hash(zeroHashes[i-1], zeroHashes[i-1]))
	}
	return zeroHashes
}
