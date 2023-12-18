package state

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/core/types"
)

const newL2BlocksCheckInterval = 200 * time.Millisecond

// NewL2BlockEventHandler represent a func that will be called by the
// state when a NewL2BlockEvent is triggered
type NewL2BlockEventHandler func(e NewL2BlockEvent)

// NewL2BlockEvent is a struct provided from the state to the NewL2BlockEventHandler
// when a new l2 block is detected with data related to this new l2 block.
type NewL2BlockEvent struct {
	Block types.Block
	Logs  []*types.Log
}

// StartToMonitorNewL2Blocks starts 2 go routines that will
// monitor new blocks and execute handlers registered to be executed
// when a new l2 block is detected. This is used by the RPC WebSocket
// filter subscription but can be used by any other component that
// needs to react to a new L2 block added to the state.
func (s *State) StartToMonitorNewL2Blocks() {
	go InfiniteSafeRun(s.monitorNewL2Blocks, "fail to monitor new l2 blocks: %v:", time.Second)
	go InfiniteSafeRun(s.handleEvents, "fail to handle events: %v", time.Second)
}

// RegisterNewL2BlockEventHandler add the provided handler to the list of handlers
// that will be triggered when a new l2 block event is triggered
func (s *State) RegisterNewL2BlockEventHandler(h NewL2BlockEventHandler) {
	log.Info("new l2 block event handler registered")
	s.newL2BlockEventHandlers = append(s.newL2BlockEventHandlers, h)
}

func (s *State) monitorNewL2Blocks() {
	waitNextCycle := func() {
		time.Sleep(newL2BlocksCheckInterval)
	}

	lastL2BlockNumber, err := s.GetLastL2BlockNumber(context.Background(), nil)
	if errors.Is(err, ErrStateNotSynchronized) {
		lastL2BlockNumber = 0
	} else if err != nil {
		log.Fatalf("failed to load the last l2 block: %v", err)
	}
	lastL2BlockNumberSeen := lastL2BlockNumber

	for {
		if len(s.newL2BlockEventHandlers) == 0 {
			waitNextCycle()
			continue
		}

		lastL2BlockNumber, err := s.GetLastL2BlockNumber(context.Background(), nil)
		if errors.Is(err, ErrStateNotSynchronized) {
			waitNextCycle()
			continue
		} else if err != nil {
			log.Errorf("failed to get last l2 block while monitoring new blocks: %v", err)
			waitNextCycle()
			continue
		}

		// not updates until now
		if lastL2BlockNumber == 0 || lastL2BlockNumberSeen >= lastL2BlockNumber {
			waitNextCycle()
			continue
		}

		fromBlockNumber := lastL2BlockNumberSeen + uint64(1)
		toBlockNumber := lastL2BlockNumber
		log.Infof("[monitorNewL2Blocks] new l2 block detected from block %v to %v", fromBlockNumber, toBlockNumber)

		for bn := fromBlockNumber; bn <= toBlockNumber; bn++ {
			block, err := s.GetL2BlockByNumber(context.Background(), bn, nil)
			if err != nil {
				log.Errorf("failed to get l2 block while monitoring new blocks: %v", err)
				break
			}
			logs, err := s.GetLogsByBlockNumber(context.Background(), bn, nil)
			if err != nil {
				log.Errorf("failed to get l2 block while monitoring new blocks: %v", err)
				break
			}

			log.Debugf("[monitorNewL2Blocks] sending NewL2BlockEvent for block %v", block.NumberU64())
			start := time.Now()
			s.newL2BlockEvents <- NewL2BlockEvent{
				Block: *block,
				Logs:  logs,
			}
			lastL2BlockNumberSeen = block.NumberU64()
			log.Infof("[monitorNewL2Blocks] NewL2BlockEvent for block %v took %v to be sent", block.NumberU64(), time.Since(start))
			log.Infof("new l2 block detected: number %v, hash %v", block.NumberU64(), block.Hash().String())
		}

		// interval to check for new l2 blocks
		waitNextCycle()
	}
}

func (s *State) handleEvents() {
	for newL2BlockEvent := range s.newL2BlockEvents {
		log.Infof("[handleEvents] new l2 block event detected for block: %v", newL2BlockEvent.Block.NumberU64())
		if len(s.newL2BlockEventHandlers) == 0 {
			continue
		}

		wg := sync.WaitGroup{}
		for _, handler := range s.newL2BlockEventHandlers {
			wg.Add(1)
			go func(h NewL2BlockEventHandler, e NewL2BlockEvent) {
				defer func() {
					wg.Done()
					if r := recover(); r != nil {
						log.Errorf("failed and recovered in NewL2BlockEventHandler: %v", r)
					}
				}()
				log.Infof("[handleEvents] triggering new l2 block event handler for block: %v", e.Block.NumberU64())
				start := time.Now()
				h(e)
				log.Infof("[handleEvents] new l2 block event handler for block %v took %v to be executed", e.Block.NumberU64(), time.Since(start))
			}(handler, newL2BlockEvent)
		}
		wg.Wait()
	}
}
