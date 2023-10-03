package pool

import "github.com/ethereum/go-ethereum/common"

// PolicyName is a named policy
type PolicyName string

const (
	// SendTx is the name of the policy that governs that an address may send transactions to pool
	SendTx PolicyName = "send_tx"
	// Deploy is the name of the policy that governs that an address may deploy a contract
	Deploy PolicyName = "deploy"
)

// Policy describes state of a named policy
type Policy struct {
	Name  PolicyName
	Allow bool
}

// Desc returns the string representation of a policy rule
func (p *Policy) Desc() string {
	if p.Allow {
		return "allow"
	}
	return "deny"
}

// Acl describes exception to a named Policy by address
type Acl struct {
	PolicyName PolicyName
	Address    common.Address
}

// IsPolicy tests if a string represents a known named Policy
func IsPolicy(name string) bool {
	for _, p := range []PolicyName{SendTx, Deploy} {
		if name == string(p) {
			return true
		}
	}
	return false
}
