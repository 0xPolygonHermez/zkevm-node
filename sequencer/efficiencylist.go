package sequencer

import (
	"fmt"
	"sort"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// efficiencyList represents a list of tx sorted by efficiency
type efficiencyList struct {
	list   map[string]*TxTracker
	sorted []*TxTracker
	mutex  sync.Mutex
}

// newEfficiencyList creates and init an efficiencyList
func newEfficiencyList() *efficiencyList {
	return &efficiencyList{
		list:   make(map[string]*TxTracker),
		sorted: []*TxTracker{},
	}
}

// add adds a tx to the efficiencyList
func (e *efficiencyList) add(tx *TxTracker) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if _, found := e.list[tx.HashStr]; !found {
		e.list[tx.HashStr] = tx
		e.addSort(tx)
		return true
	}
	return false
}

// delete deletes the tx from the efficiencyList
func (e *efficiencyList) delete(tx *TxTracker) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if tx, found := e.list[tx.HashStr]; found {
		sLen := len(e.sorted)
		i := sort.Search(sLen, func(i int) bool {
			return e.isGreaterThan(tx, e.list[e.sorted[i].HashStr])
		})

		if (e.sorted[i].HashStr != tx.HashStr) || i == sLen {
			log.Errorf("Error deleting tx from efficiencyList: %s", tx.HashStr)
			return false
		}

		delete(e.list, tx.HashStr)

		copy(e.sorted[i:], e.sorted[i+1:])
		e.sorted[sLen-1] = nil
		e.sorted = e.sorted[:sLen-1]

		return true
	}
	return false
}

// getByIndex retrieves the tx at the i position in the sorted EfficiencyList
func (e *efficiencyList) getByIndex(i int) *TxTracker {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	tx := e.sorted[i]

	return tx
}

// len returns the length of the EfficiencyList
func (e *efficiencyList) len() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	l := len(e.sorted)

	return l
}

// print prints the contents of the EfficiencyList
func (e *efficiencyList) Print() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	fmt.Println("Len: ", len(e.sorted))
	for _, txi := range e.sorted {
		fmt.Printf("Hash=%s, efficiency=%f\n", txi.HashStr, txi.Efficiency)
	}
}

// addSort adds the tx to the EfficiencyList in a sorted way
func (e *efficiencyList) addSort(tx *TxTracker) {
	i := sort.Search(len(e.sorted), func(i int) bool {
		return e.isGreaterThan(tx, e.list[e.sorted[i].HashStr])
	})

	e.sorted = append(e.sorted, nil)
	copy(e.sorted[i+1:], e.sorted[i:])
	e.sorted[i] = tx
	//log.Infof("Added tx(%s) to efficiencyList. With efficiency(%f) at index(%d) from total(%d)", tx.HashStr, tx.Efficiency, i, len(e.sorted))
}

// isGreaterThan returns true if the tx1 has best efficiency than tx2
func (e *efficiencyList) isGreaterThan(tx1 *TxTracker, tx2 *TxTracker) bool {
	if tx1.Efficiency > tx2.Efficiency {
		return true
	} else if tx1.Efficiency == tx2.Efficiency {
		return tx1.HashStr >= tx2.HashStr
	} else {
		return false
	}
}

// GetSorted returns the sorted list of tx
func (e *efficiencyList) GetSorted() []*TxTracker {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.sorted
}
