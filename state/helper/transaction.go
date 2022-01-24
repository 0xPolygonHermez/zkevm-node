package helper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func GetSender(tx *types.Transaction) (common.Address, error) {
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}
	return sender, nil
}
