package l1_parallel_sync

import (
	"errors"
	"fmt"
)

const (
	latestBlockNumber  uint64 = ^uint64(0)
	invalidBlockNumber uint64 = uint64(0)
)

var (
	errBlockRangeInvalidIsNil   = errors.New("block Range Invalid: block range is nil")
	errBlockRangeInvalidIsZero  = errors.New("block Range Invalid: Invalid: from or to are 0")
	errBlockRangeInvalidIsWrong = errors.New("block Range Invalid: fromBlock is greater than toBlock")
)

type blockRange struct {
	fromBlock uint64
	toBlock   uint64
}

func blockNumberToString(b uint64) string {
	if b == latestBlockNumber {
		return "latest"
	}
	if b == invalidBlockNumber {
		return "invalid"
	}
	return fmt.Sprintf("%d", b)
}

func (b *blockRange) String() string {
	return fmt.Sprintf("[%s, %s]", blockNumberToString(b.fromBlock), blockNumberToString(b.toBlock))
}

func (b *blockRange) len() uint64 {
	if b.toBlock == latestBlockNumber || b.fromBlock == latestBlockNumber {
		return 0
	}
	return b.toBlock - b.fromBlock + 1
}

func (b *blockRange) isValid() error {
	if b == nil {
		return errBlockRangeInvalidIsNil
	}
	if b.fromBlock == invalidBlockNumber || b.toBlock == invalidBlockNumber {
		return errBlockRangeInvalidIsZero
	}
	if b.fromBlock > b.toBlock {
		return errBlockRangeInvalidIsWrong
	}
	return nil
}

func (b *blockRange) overlaps(br blockRange) bool {
	return b.fromBlock <= br.toBlock && br.fromBlock <= b.toBlock
}
