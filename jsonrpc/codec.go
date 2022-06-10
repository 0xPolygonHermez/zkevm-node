package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hermeznetwork/hermez-core/encoding"
)

const (
	// PendingBlockNumber represents the pending block number
	PendingBlockNumber = BlockNumber(-3)
	// LatestBlockNumber represents the latest block number
	LatestBlockNumber = BlockNumber(-2)
	// EarliestBlockNumber represents the earliest block number
	EarliestBlockNumber = BlockNumber(-1)
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

func (b *BlockNumber) getNumericBlockNumber(ctx context.Context, s stateInterface) (uint64, error) {
	if b == nil {
		return 0, nil
	}

	bValue := *b
	switch bValue {
	case LatestBlockNumber, PendingBlockNumber:
		lastBatchNumber, err := s.GetLastBatchNumber(ctx, "")
		if err != nil {
			return 0, err
		}

		return lastBatchNumber, nil

	case EarliestBlockNumber:
		return 0, nil

	default:
		if bValue < 0 {
			return 0, fmt.Errorf("invalid block number: %v", bValue)
		}
		return uint64(bValue), nil
	}
}

func stringToBlockNumber(str string) (BlockNumber, error) {
	str = strings.Trim(str, "\"")
	switch str {
	case "earliest":
		return EarliestBlockNumber, nil
	case "pending":
		return PendingBlockNumber, nil
	case "latest", "":
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
