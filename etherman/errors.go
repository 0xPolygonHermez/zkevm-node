package etherman

import "errors"

var (
	//ErrGasRequiredExceedsAllowance gas required exceeds the allowance
	ErrGasRequiredExceedsAllowance = errors.New("gas required exceeds allowance")
	//ErrContentLengthTooLarge content length is too large
	ErrContentLengthTooLarge = errors.New("content length too large")
	//ErrTimestampMustBeInsideRange Timestamp must be inside range
	ErrTimestampMustBeInsideRange = errors.New("Timestamp must be inside range")
	//ErrInsufficientAllowance insufficient allowance
	ErrInsufficientAllowance = errors.New("insufficient allowance")
	//ErrBothGasPriceAndMaxFeeGasAreSpecified both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified
	ErrBothGasPriceAndMaxFeeGasAreSpecified = errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	//ErrMaxFeeGasAreSpecifiedButLondonNotActive maxFeePerGas or maxPriorityFeePerGas specified but london is not active yet
	ErrMaxFeeGasAreSpecifiedButLondonNotActive = errors.New("maxFeePerGas or maxPriorityFeePerGas specified but london is not active yet")
	//ErrNoSigner no signer to authorize the transaction with
	ErrNoSigner = errors.New("no signer to authorize the transaction with")

	errorsCache = map[string]error{
		ErrGasRequiredExceedsAllowance.Error():             ErrGasRequiredExceedsAllowance,
		ErrContentLengthTooLarge.Error():                   ErrContentLengthTooLarge,
		ErrTimestampMustBeInsideRange.Error():              ErrTimestampMustBeInsideRange,
		ErrInsufficientAllowance.Error():                   ErrInsufficientAllowance,
		ErrBothGasPriceAndMaxFeeGasAreSpecified.Error():    ErrBothGasPriceAndMaxFeeGasAreSpecified,
		ErrMaxFeeGasAreSpecifiedButLondonNotActive.Error(): ErrMaxFeeGasAreSpecifiedButLondonNotActive,
		ErrNoSigner.Error():                                ErrNoSigner,
	}
)

func tryParseError(err error) (error, bool) {
	parsedError, exists := errorsCache[err.Error()]
	return parsedError, exists
}
