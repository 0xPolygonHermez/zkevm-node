package merkletree

// leafType specifies type of the leaf
type leafType uint8

const (
	// leafTypeBalance specifies that leaf stores Balance
	leafTypeBalance leafType = 0
	// leafTypeNonce specifies that leaf stores Nonce
	leafTypeNonce leafType = 1
	// leafTypeCode specifies that leaf stores Code
	leafTypeCode leafType = 2
	// leafTypeStorage specifies that leaf stores Storage Value
	leafTypeStorage leafType = 3
)
