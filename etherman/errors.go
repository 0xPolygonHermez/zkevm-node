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

	errorsFromPoEContract = map[string]error{
		ErrGasRequiredExceedsAllowance.Error(): ErrGasRequiredExceedsAllowance,
		ErrContentLengthTooLarge.Error():       ErrContentLengthTooLarge,
		ErrTimestampMustBeInsideRange.Error():  ErrTimestampMustBeInsideRange,
		ErrInsufficientAllowance.Error():       ErrInsufficientAllowance,
	}
)

func tryParseContractPoEError(err error) (error, bool) {
	parsedError, exists := errorsFromPoEContract[err.Error()]
	return parsedError, exists
}
