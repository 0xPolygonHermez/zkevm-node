package tree

// LeafType specifies type of the leaf
type LeafType uint8

const (
	// LeafTypeBalance specifies that leaf stores Balance
	LeafTypeBalance LeafType = 0
	// LeafTypeNonce specifies that leaf stores Nonce
	LeafTypeNonce LeafType = 1
	// LeafTypeCode specifies that leaf stores Code
	LeafTypeCode LeafType = 2
	// LeafTypeStorage specifies that leaf stores Storage Value
	LeafTypeStorage LeafType = 3
)
