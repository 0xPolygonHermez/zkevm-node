package jsonrpc

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/jackc/pgx/v4"
)

const (
	// PendingBlockNumber represents the pending block number
	PendingBlockNumber = BlockNumber(-3)
	// LatestBlockNumber represents the latest block number
	LatestBlockNumber = BlockNumber(-2)
	// EarliestBlockNumber represents the earliest block number
	EarliestBlockNumber = BlockNumber(-1)

	// Earliest contains the string to represent the earliest block known.
	Earliest = "earliest"
	// Latest contains the string to represent the latest block known.
	Latest = "latest"
	// Pending contains the string to represent pending blocks.
	Pending = "pending"
)

// Request is a jsonrpc request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response is a jsonrpc  success response
type Response struct {
	JSONRPC string
	ID      interface{}
	Result  json.RawMessage
	Error   *ErrorObject
}

// ErrorObject is a jsonrpc error
type ErrorObject struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewResponse returns Success/Error response object
func NewResponse(req Request, reply *[]byte, err rpcError) Response {
	var result json.RawMessage
	if reply != nil {
		result = *reply
	}

	var errorObj *ErrorObject
	if err != nil {
		errorObj = &ErrorObject{err.ErrorCode(), err.Error(), nil}
	}

	return Response{
		JSONRPC: req.JSONRPC,
		ID:      req.ID,
		Result:  result,
		Error:   errorObj,
	}
}

// MarshalJSON customizes the JSON representation of the response.
func (r Response) MarshalJSON() ([]byte, error) {
	if r.Error != nil {
		return json.Marshal(struct {
			JSONRPC string       `json:"jsonrpc"`
			ID      interface{}  `json:"id"`
			Error   *ErrorObject `json:"error"`
		}{
			JSONRPC: r.JSONRPC,
			ID:      r.ID,
			Error:   r.Error,
		})
	}

	return json.Marshal(struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      interface{}     `json:"id"`
		Result  json.RawMessage `json:"result"`
	}{
		JSONRPC: r.JSONRPC,
		ID:      r.ID,
		Result:  r.Result,
	})
}

// BlockNumber is the number of a ethereum block
type BlockNumber int64

// UnmarshalJSON automatically decodes the user input for the block number, when a JSON RPC method is called
func (b *BlockNumber) UnmarshalJSON(buffer []byte) error {
	num, err := stringToBlockNumber(string(buffer))
	if err != nil {
		return err
	}
	*b = num
	return nil
}

func (b *BlockNumber) getNumericBlockNumber(ctx context.Context, s stateInterface, dbTx pgx.Tx) (uint64, rpcError) {
	bValue := LatestBlockNumber
	if b != nil {
		bValue = *b
	}

	switch bValue {
	case LatestBlockNumber, PendingBlockNumber:
		lastBlockNumber, err := s.GetLastL2BlockNumber(ctx, dbTx)
		if err != nil {
			return 0, newRPCError(defaultErrorCode, "failed to get the last block number from state")
		}

		return lastBlockNumber, nil

	case EarliestBlockNumber:
		return 0, nil

	default:
		if bValue < 0 {
			return 0, newRPCError(invalidParamsErrorCode, "invalid block number: %v", bValue)
		}
		return uint64(bValue), nil
	}
}

func stringToBlockNumber(str string) (BlockNumber, error) {
	str = strings.Trim(str, "\"")
	switch str {
	case Earliest:
		return EarliestBlockNumber, nil
	case Pending:
		return PendingBlockNumber, nil
	case Latest, "":
		return LatestBlockNumber, nil
	}

	n, err := encoding.DecodeUint64orHex(&str)
	if err != nil {
		return 0, err
	}
	return BlockNumber(n), nil
}

// Index of a item
type Index int64

// UnmarshalJSON automatically decodes the user input for the block number, when a JSON RPC method is called
func (i *Index) UnmarshalJSON(buffer []byte) error {
	str := strings.Trim(string(buffer), "\"")
	n, err := encoding.DecodeUint64orHex(&str)
	if err != nil {
		return err
	}
	*i = Index(n)
	return nil
}
