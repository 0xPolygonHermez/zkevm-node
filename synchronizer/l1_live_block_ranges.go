package synchronizer

import (
	"errors"
	"fmt"
)

type liveBlockRangeItem[T any] struct {
	blockRange blockRange
	tag        T
}

type liveBlockRangesGeneric[T any] struct {
	ranges []liveBlockRangeItem[T]
}

type liveBlockRanges = liveBlockRangesGeneric[int]

func (l *liveBlockRangesGeneric[T]) String() string {
	res := l.toStringBrief() + "["
	for _, r := range l.ranges {
		res += fmt.Sprintf("%s ,", r.blockRange.String())
	}
	return res + "]"
}

func (l *liveBlockRangesGeneric[T]) toStringBrief() string {
	return fmt.Sprintf("len(ranges): %v", len(l.ranges))
}

const (
	invalidBlockNumber = uint64(0)
)

var (
	errBlockRangeInvalidIsNil   = errors.New("block Range Invalid: block range is nil")
	errBlockRangeInvalidIsZero  = errors.New("block Range Invalid: Invalid: from or to are 0")
	errBlockRangeInvalidIsWrong = errors.New("block Range Invalid: fromBlock is greater than toBlock")
	errBlockRangeInvalidOverlap = errors.New("block Range Invalid: block range overlaps")
	errBlockRangeNotFound       = errors.New("block Range not found")
	errBlockRangeIsEmpty        = errors.New("block Range is empty")
)

func newLiveBlockRanges() liveBlockRanges {
	return liveBlockRanges{}
}

func newLiveBlockRangesWithTag[T any]() liveBlockRangesGeneric[T] {
	return liveBlockRangesGeneric[T]{}
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

func (l *liveBlockRangesGeneric[T]) addBlockRange(br blockRange) error {
	var zeroValue T
	return l.addBlockRangeWithTag(br, zeroValue)
}

func (l *liveBlockRangesGeneric[T]) addBlockRangeWithTag(br blockRange, tag T) error {
	if err := br.isValid(); err != nil {
		return err
	}
	if l.overlaps(br) {
		return errBlockRangeInvalidOverlap
	}
	l.ranges = append(l.ranges, liveBlockRangeItem[T]{blockRange: br, tag: tag})
	return nil
}

func (l *liveBlockRangesGeneric[T]) filterBlockRangesByTag(filter func(blockRange, T) bool) []blockRange {
	result := make([]blockRange, 0)
	for _, r := range l.ranges {
		if filter(r.blockRange, r.tag) {
			result = append(result, r.blockRange)
		}
	}
	return result
}

func (l *liveBlockRangesGeneric[T]) removeBlockRange(br blockRange) error {
	for i, r := range l.ranges {
		if r.blockRange == br {
			l.ranges = append(l.ranges[:i], l.ranges[i+1:]...)
			return nil
		}
	}
	return errBlockRangeNotFound
}

func (l *liveBlockRangesGeneric[T]) getFirstBlockRange() (blockRange, error) {
	if l.len() == 0 {
		return blockRange{}, errBlockRangeIsEmpty
	}
	return l.ranges[0].blockRange, nil
}

func (l *liveBlockRangesGeneric[T]) getTagByBlockRange(br blockRange) (T, error) {
	for _, r := range l.ranges {
		if r.blockRange == br {
			return r.tag, nil
		}
	}
	var zeroValue T
	return zeroValue, errBlockRangeNotFound
}

func (l *liveBlockRangesGeneric[T]) getSuperBlockRange() *blockRange {
	fromBlock := invalidBlockNumber
	toBlock := invalidBlockNumber
	for i, r := range l.ranges {
		if i == 0 {
			toBlock = r.blockRange.toBlock
			fromBlock = r.blockRange.fromBlock
		}
		if r.blockRange.toBlock > toBlock {
			toBlock = r.blockRange.toBlock
		}
		if r.blockRange.fromBlock < fromBlock {
			fromBlock = r.blockRange.fromBlock
		}
	}
	res := blockRange{fromBlock, toBlock}
	if res.isValid() == nil {
		return &res
	}
	return nil
}

func (l *liveBlockRangesGeneric[T]) len() int {
	return len(l.ranges)
}

func (l *liveBlockRangesGeneric[T]) overlaps(br blockRange) bool {
	for _, r := range l.ranges {
		if r.blockRange.overlaps(br) {
			return true
		}
	}
	return false
}
