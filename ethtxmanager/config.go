package ethtxmanager

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config is configuration for ethereum transaction manager
type Config struct {
	// FrequencyToMonitorTxs frequency of the resending failed txs
	FrequencyToMonitorTxs types.Duration `mapstructure:"FrequencyToMonitorTxs"`
	// WaitTxToBeMined time to wait after transaction was sent to the ethereum
	WaitTxToBeMined types.Duration `mapstructure:"WaitTxToBeMined"`

	// PrivateKeys defines all the key store files that are going
	// to be read in order to provide the private keys to sign the L1 txs
	PrivateKeys []types.KeystoreFileConfig `mapstructure:"PrivateKeys"`

	// ForcedGas is the amount of gas to be forced in case of gas estimation error
	ForcedGas uint64 `mapstructure:"ForcedGas"`

	// GasPriceMarginFactor is used to multiply the suggested gas price provided by the network
	// in order to allow a different gas price to be set for all the transactions and make it
	// easier to have the txs prioritized in the pool, the default value is 1.
	//
	// ex:
	// suggested gas price: 100
	// GasPriceMarginFactor: 1
	// gas price = 100
	//
	// suggested gas price: 100
	// GasPriceMarginFactor: 1.1
	// gas price = 110
	GasPriceMarginFactor float64 `mapstructure:"GasPriceMarginFactor"`

	// MaxGasPriceLimit helps to avoid transactions to be sent over a specified
	// gas price amount, the default value is 0, which means no limit.
	// If the gas price provided by the network and adjusted by the GasPriceMarginFactor
	// is greater than this configuration, the transaction will have its gas price set to
	// the value configured in this config as the limit.
	//
	// ex:
	//
	// suggested gas price: 100
	// gas price margin factor: 20%
	// max gas price limit: 150
	// tx gas price = 120
	//
	// suggested gas price: 100
	// gas price margin factor: 20%
	// max gas price limit: 110
	// tx gas price = 110
	MaxGasPriceLimit uint64 `mapstructure:"MaxGasPriceLimit"`
}
