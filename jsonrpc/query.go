package jsonrpc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Filter represents a filter.
type Filter struct {
	ID         uint64
	Type       string
	Parameters string
	LastPoll   time.Time
}

// LogFilterRequest represents a log filter request.
type LogFilterRequest struct {
	BlockHash *common.Hash  `json:"blockHash,omitempty"`
	FromBlock string        `json:"fromBlock,omitempty"`
	ToBlock   string        `json:"toBlock,omitempty"`
	Address   interface{}   `json:"address,omitempty"`
	Topics    []interface{} `json:"topics,omitempty"`
}

// LogFilter is a filter for logs
type LogFilter struct {
	BlockHash *common.Hash
	FromBlock BlockNumber
	ToBlock   BlockNumber
	Addresses []common.Address
	Topics    [][]common.Hash
	Since     *time.Time
}

// addTopic adds specific topics to the log filter topics
func (f *LogFilter) addTopic(topics ...string) error {
	if f.Topics == nil {
		f.Topics = [][]common.Hash{}
	}

	topicsHashes := []common.Hash{}

	for _, topic := range topics {
		topicHash := common.Hash{}
		if err := topicHash.UnmarshalText([]byte(topic)); err != nil {
			return err
		}

		topicsHashes = append(topicsHashes, topicHash)
	}

	f.Topics = append(f.Topics, topicsHashes)

	return nil
}

// addAddress Adds the address to the log filter
func (f *LogFilter) addAddress(raw string) error {
	if f.Addresses == nil {
		f.Addresses = []common.Address{}
	}

	addr := common.Address{}

	if err := addr.UnmarshalText([]byte(raw)); err != nil {
		return err
	}

	f.Addresses = append(f.Addresses, addr)

	return nil
}

// MarshalJSON allows to customize the JSON representation.
func (f *LogFilter) MarshalJSON() ([]byte, error) {
	var obj LogFilterRequest

	obj.BlockHash = f.BlockHash

	if f.FromBlock == LatestBlockNumber {
		obj.FromBlock = ""
	} else {
		obj.FromBlock = hex.EncodeUint64(uint64(f.FromBlock))
	}

	if f.ToBlock == LatestBlockNumber {
		obj.ToBlock = ""
	} else {
		obj.ToBlock = hex.EncodeUint64(uint64(f.ToBlock))
	}

	if f.Addresses != nil {
		if len(f.Addresses) == 1 {
			obj.Address = f.Addresses[0].Hex()
		} else {
			obj.Address = f.Addresses
		}
	}

	obj.Topics = make([]interface{}, 0, len(f.Topics))
	for _, topic := range f.Topics {
		if len(topic) == 0 {
			obj.Topics = append(obj.Topics, nil)
		} else if len(topic) == 1 {
			obj.Topics = append(obj.Topics, topic[0])
		} else {
			obj.Topics = append(obj.Topics, topic)
		}
	}

	return json.Marshal(obj)
}

// UnmarshalJSON decodes a json object
func (f *LogFilter) UnmarshalJSON(data []byte) error {
	var obj LogFilterRequest

	err := json.Unmarshal(data, &obj)

	if err != nil {
		return err
	}

	f.BlockHash = obj.BlockHash

	if obj.FromBlock == "" {
		f.FromBlock = LatestBlockNumber
	} else {
		if f.FromBlock, err = stringToBlockNumber(obj.FromBlock); err != nil {
			return err
		}
	}

	if obj.ToBlock == "" {
		f.ToBlock = LatestBlockNumber
	} else {
		if f.ToBlock, err = stringToBlockNumber(obj.ToBlock); err != nil {
			return err
		}
	}

	if obj.Address != nil {
		// decode address, either "" or [""]
		switch raw := obj.Address.(type) {
		case string:
			// ""
			if err := f.addAddress(raw); err != nil {
				return err
			}

		case []interface{}:
			// ["", ""]
			for _, addr := range raw {
				if item, ok := addr.(string); ok {
					if err := f.addAddress(item); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("address expected")
				}
			}

		default:
			return fmt.Errorf("failed to decode address. Expected either '' or ['', '']")
		}
	}

	if obj.Topics != nil {
		// decode topics, either "" or ["", ""] or null
		for _, item := range obj.Topics {
			switch raw := item.(type) {
			case string:
				// ""
				if err := f.addTopic(raw); err != nil {
					return err
				}

			case []interface{}:
				// ["", ""]
				res := []string{}

				for _, i := range raw {
					if item, ok := i.(string); ok {
						res = append(res, item)
					} else {
						return fmt.Errorf("hash expected")
					}
				}

				if err := f.addTopic(res...); err != nil {
					return err
				}

			case nil:
				// null
				if err := f.addTopic(); err != nil {
					return err
				}

			default:
				return fmt.Errorf("failed to decode topics. Expected '' or [''] or null")
			}
		}
	}

	// decode topics
	return nil
}

// Match returns whether the receipt includes topics for this filter
func (f *LogFilter) Match(log *types.Log) bool {
	// check addresses
	if len(f.Addresses) > 0 {
		match := false

		for _, addr := range f.Addresses {
			if addr == log.Address {
				match = true
			}
		}

		if !match {
			return false
		}
	}
	// check topics
	if len(f.Topics) > len(log.Topics) {
		return false
	}

	for i, sub := range f.Topics {
		match := len(sub) == 0

		for _, topic := range sub {
			if log.Topics[i] == topic {
				match = true

				break
			}
		}

		if !match {
			return false
		}
	}

	return true
}
