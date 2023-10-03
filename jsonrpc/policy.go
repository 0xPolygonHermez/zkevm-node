package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

func checkPolicy(ctx context.Context, p types.PoolInterface, input string) error {
	tx, err := hexToTx(input)
	if err != nil {
		// ignore it, let the later processing reject
		return nil
	}

	// if the tx is signed, check the from address. If there is no from address, the tx is not rejected as it
	// will get rejected later. This maintains backward compatibility with RPC expectations. TODO: verify this is ok behavior
	var from common.Address
	if from, err = state.GetSender(*tx); err != nil {
		// if not signed, then skip check, it fails later on its own
		return nil
	}

	switch resolvePolicy(tx) {
	case pool.SendTx:
		var allow bool
		if allow, err = p.CheckPolicy(ctx, pool.SendTx, *tx.To()); err != nil {
			return err
		}
		if !allow {
			return pool.ErrContractDisallowedSendTx
		}
		if allow, err = p.CheckPolicy(ctx, pool.SendTx, from); err != nil {
			return err
		}
		if !allow {
			return pool.ErrSenderDisallowedSendTx
		}
	case pool.Deploy:
		var allow bool
		// check that sender may deploy contracts
		if allow, err = p.CheckPolicy(ctx, pool.Deploy, from); err != nil {
			return err
		}
		if !allow {
			return pool.ErrSenderDisallowedDeploy
		}
	}
	return nil
}

func resolvePolicy(tx *ethTypes.Transaction) pool.PolicyName {
	if tx.To() == nil || tx.To().Hex() == common.HexToAddress("0x0").Hex() {
		return pool.Deploy
	}
	return pool.SendTx
}
