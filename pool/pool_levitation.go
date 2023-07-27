package pool

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
)

// VerifyTx withouyt adding it to the pool. It will go to Levitation chain PendingQueue contract after verification
func (p *Pool) VerifyTx(ctx context.Context, tx types.Transaction, ip string) error {
	poolTx := NewTransaction(tx, ip, false)
	return p.validateTx(ctx, *poolTx)
}
