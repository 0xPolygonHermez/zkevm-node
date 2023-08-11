package synchronizer

import (
	"errors"
	"fmt"
)

type liveBlockRangeItem struct {
	blockRange blockRange
}

type liveBlockRanges struct {
	ranges []liveBlockRangeItem
}

func (l *liveBlockRanges) toString() string {
	res := l.toStringBrief() + "["
	for _, r := range l.ranges {
		res += fmt.Sprintf("%s ,", r.blockRange.toString())
	}
	return res + "]"
}

func (l *liveBlockRanges) toStringBrief() string {
	return fmt.Sprintf("len(ranges): %v", len(l.ranges))
}

const (
	invalidBlockNumber = uint64(0)
)

const (
	errBlockRangeInvalidIsNil   = "block Range Invalid: block range is nil"
	errBlockRangeInvalidIsZero  = "block Range Invalid: Invalid: from or to are 0"
	errBlockRangeInvalidIsWrong = "block Range Invalid: fromBlock is greater than toBlock"
	errBlockRangeInvalidOverlap = "block Range Invalid: block range overlaps"
	errBlockRangeNotFound       = "block Range not found"
	errBlockRangeIsEmpty        = "block Range is empty"
)

func newLiveBlockRanges() liveBlockRanges {
	return liveBlockRanges{}
}

func (b *blockRange) isValid() error {
	if b == nil {
		return errors.New(errBlockRangeInvalidIsNil)
	}
	if b.fromBlock == invalidBlockNumber || b.toBlock == invalidBlockNumber {
		return errors.New(errBlockRangeInvalidIsZero)
	}
	if b.fromBlock > b.toBlock {
		return errors.New(errBlockRangeInvalidIsWrong + b.toString())
	}
	return nil
}

func (b *blockRange) overlaps(br blockRange) bool {
	return b.fromBlock <= br.toBlock && br.fromBlock <= b.toBlock
}

func (l *liveBlockRanges) addBlockRange(br blockRange) error {
	if err := br.isValid(); err != nil {
		return err
	}
	if l.overlaps(br) {
		return errors.New(errBlockRangeInvalidOverlap)
	}
	l.ranges = append(l.ranges, liveBlockRangeItem{br})
	return nil
}

func (l *liveBlockRanges) removeBlockRange(br blockRange) error {
	for i, r := range l.ranges {
		if r.blockRange == br {
			l.ranges = append(l.ranges[:i], l.ranges[i+1:]...)
			return nil
		}
	}
	return errors.New(errBlockRangeNotFound)
}

func (l *liveBlockRanges) getFirstBlockRange() (blockRange, error) {
	if l.len() == 0 {
		return blockRange{}, errors.New(errBlockRangeIsEmpty)
	}
	return l.ranges[0].blockRange, nil
}

func (l *liveBlockRanges) getSuperBlockRange() *blockRange {
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

func (l *liveBlockRanges) len() int {
	return len(l.ranges)
}

func (l *liveBlockRanges) overlaps(br blockRange) bool {
	for _, r := range l.ranges {
		if r.blockRange.overlaps(br) {
			return true
		}
	}
	return false
}
