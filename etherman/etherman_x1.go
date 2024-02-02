package etherman

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// EstimateGasSequenceBatchesX1 estimates gas for sending batches
func (etherMan *Client) EstimateGasSequenceBatchesX1(sender common.Address, sequences []ethmanTypes.Sequence, l2Coinbase common.Address, committeeSignaturesAndAddrs []byte) (*types.Transaction, error) {
	opts, err := etherMan.generateMockAuth(sender)
	if err == ErrNotFound {
		return nil, ErrPrivateKeyNotFound
	}
	opts.NoSend = true

	tx, err := etherMan.sequenceBatchesX1(opts, sequences, l2Coinbase, committeeSignaturesAndAddrs)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// BuildSequenceBatchesTxDataX1 builds a []bytes to be sent to the PoE SC method SequenceBatches.
func (etherMan *Client) BuildSequenceBatchesTxDataX1(sender common.Address, sequences []ethmanTypes.Sequence, l2Coinbase common.Address, committeeSignaturesAndAddrs []byte) (to *common.Address, data []byte, err error) {
	opts, err := etherMan.generateRandomAuth()
	if err == ErrNotFound {
		return nil, nil, fmt.Errorf("failed to build sequence batches, err: %w", ErrPrivateKeyNotFound)
	}
	opts.NoSend = true
	// force nonce, gas limit and gas price to avoid querying it from the chain
	opts.Nonce = big.NewInt(1)
	opts.GasLimit = uint64(1)
	opts.GasPrice = big.NewInt(1)

	tx, err := etherMan.sequenceBatchesX1(opts, sequences, l2Coinbase, committeeSignaturesAndAddrs)
	if err != nil {
		return nil, nil, err
	}

	return tx.To(), tx.Data(), nil
}

func (etherMan *Client) sequenceBatchesX1(opts bind.TransactOpts, sequences []ethmanTypes.Sequence, l2Coinbase common.Address, committeeSignaturesAndAddrs []byte) (*types.Transaction, error) {
	var batches []polygonzkevm.PolygonZkEVMBatchData

	var tx *types.Transaction
	var err error
	if len(committeeSignaturesAndAddrs) > 0 {
		for _, seq := range sequences {
			batch := polygonzkevm.PolygonZkEVMBatchData{
				TransactionsHash:   crypto.Keccak256Hash(seq.BatchL2Data),
				GlobalExitRoot:     seq.GlobalExitRoot,
				Timestamp:          uint64(seq.Timestamp),
				MinForcedTimestamp: uint64(seq.ForcedBatchTimestamp),
			}

			batches = append(batches, batch)
		}

		log.Infof("Sequence batches with validium.")
		tx, err = etherMan.ZkEVM.SequenceBatches(&opts, batches, l2Coinbase, committeeSignaturesAndAddrs)
	} else {
		for _, seq := range sequences {
			batch := polygonzkevm.PolygonZkEVMBatchData{
				Transactions:       seq.BatchL2Data,
				GlobalExitRoot:     seq.GlobalExitRoot,
				Timestamp:          uint64(seq.Timestamp),
				MinForcedTimestamp: uint64(seq.ForcedBatchTimestamp),
			}

			batches = append(batches, batch)
		}

		log.Infof("Sequence batches with rollup.")
		tx, err = etherMan.ZkEVM.SequenceBatches(&opts, batches, l2Coinbase, nil)
	}

	if err != nil {
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
		err = fmt.Errorf(
			"error sequencing batches: %w, committeeSignaturesAndAddrs %s",
			err, common.Bytes2Hex(committeeSignaturesAndAddrs),
		)
	}

	return tx, err
}

// LoadAuthFromKeyStoreX1 loads an authorization from a key store file
func (etherMan *Client) LoadAuthFromKeyStoreX1(path, password string) (*bind.TransactOpts, *ecdsa.PrivateKey, error) {
	auth, pk, err := newAuthFromKeystoreX1(path, password, etherMan.l1Cfg.L1ChainID)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("loaded authorization for address: %v", auth.From.String())
	etherMan.auth[auth.From] = auth
	return &auth, pk, nil
}

// newAuthFromKeystoreX1 an authorization instance from a keystore file
func newAuthFromKeystoreX1(path, password string, chainID uint64) (bind.TransactOpts, *ecdsa.PrivateKey, error) {
	log.Infof("reading key from: %v", path)
	key, err := newKeyFromKeystore(path, password)
	if err != nil {
		return bind.TransactOpts{}, nil, err
	}
	if key == nil {
		return bind.TransactOpts{}, nil, nil
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, new(big.Int).SetUint64(chainID))
	if err != nil {
		return bind.TransactOpts{}, nil, err
	}
	return *auth, key.PrivateKey, nil
}

// generateMockAuth generates an authorization instance from a
// randomly generated private key to be used to estimate gas for PoE
// operations NOT restricted to the Trusted Sequencer
func (etherMan *Client) generateMockAuth(sender common.Address) (bind.TransactOpts, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return bind.TransactOpts{}, errors.New("failed to generate a private key to estimate L1 txs")
	}
	chainID := big.NewInt(0).SetUint64(etherMan.l1Cfg.L1ChainID)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return bind.TransactOpts{}, errors.New("failed to generate a fake authorization to estimate L1 txs")
	}

	auth.From = sender
	auth.Signer = func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		chainID := big.NewInt(0).SetUint64(etherMan.l1Cfg.L1ChainID)
		signer := types.LatestSignerForChainID(chainID)
		if err != nil {
			return nil, err
		}
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), privateKey)
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}
	return *auth, nil
}
