package sequencer

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"
)

// randomFloat64 is a shortcut for generating a random float between 0 and 1 using crypto/rand.
func randomFloat64() float64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1<<53))
	if err != nil {
		panic(err)
	}
	return float64(nBig.Int64()) / (1 << 53)
}

func TestEfficiencyListSort(t *testing.T) {
	el := newEfficiencyList()
	nItems := 100

	for i := 0; i < nItems; i++ {
		el.add(&TxTracker{HashStr: fmt.Sprintf("0x%d", i), Efficiency: randomFloat64()})
	}

	for i := 0; i < nItems-1; i++ {
		if !(el.getByIndex(i).Efficiency > el.getByIndex(i+1).Efficiency) {
			t.Fatalf("Sort error. [%d].Efficiency(%f) < [%d].Efficiency(%f)", i, el.getByIndex(i).Efficiency, i+1, el.getByIndex(i+1).Efficiency)
		}
	}

	// el.print()

	if el.len() != nItems {
		t.Fatalf("Length error. Length %d. Expected %d", el.len(), nItems)
	}
}

func TestEfficiencyListDelete(t *testing.T) {
	el := newEfficiencyList()

	el.add(&TxTracker{HashStr: "0x01", Efficiency: 1})
	el.add(&TxTracker{HashStr: "0x02", Efficiency: 2})
	el.add(&TxTracker{HashStr: "0x03", Efficiency: 2})
	el.add(&TxTracker{HashStr: "0x04", Efficiency: 3})
	el.add(&TxTracker{HashStr: "0x05", Efficiency: 10})
	el.add(&TxTracker{HashStr: "0x06", Efficiency: 1.5})
	el.add(&TxTracker{HashStr: "0x07", Efficiency: 1.5})

	deltxs := []string{"0x03", "0x07", "0x01", "0x05"}

	for _, deltx := range deltxs {
		count := el.len()
		el.delete(&TxTracker{HashStr: deltx})

		for i := 0; i < el.len(); i++ {
			if el.getByIndex(i).HashStr == deltx {
				t.Fatalf("Delete error. %s tx was not deleted", deltx)
			}
		}

		if el.len() != count-1 {
			t.Fatalf("Length error. Length %d. Expected %d", el.len(), count)
		}
	}
}

func TestEfficiencyListBench(t *testing.T) {
	el := newEfficiencyList()

	start := time.Now()
	for i := 0; i < 10000; i++ {
		el.add(&TxTracker{HashStr: fmt.Sprintf("0x%d", i), Efficiency: randomFloat64()})
	}
	elapsed := time.Since(start)
	t.Logf("EfficiencyList adding 10000 items took %s", elapsed)

	start = time.Now()
	el.add(&TxTracker{HashStr: fmt.Sprintf("0x%d", 10001), Efficiency: randomFloat64()})
	elapsed = time.Since(start)
	t.Logf("EfficiencyList adding the 10001 item (efficiency=random) took %s", elapsed)

	start = time.Now()
	el.add(&TxTracker{HashStr: fmt.Sprintf("0x%d", 10002), Efficiency: 0})
	elapsed = time.Since(start)
	t.Logf("EfficiencyList adding the 10002 item (efficiency=0) took %s", elapsed)

	start = time.Now()
	el.add(&TxTracker{HashStr: fmt.Sprintf("0x%d", 10003), Efficiency: 1000})
	elapsed = time.Since(start)
	t.Logf("EfficiencyList adding the 10003 item (efficiency=1000) took %s", elapsed)
}
