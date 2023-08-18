package sequencer

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"
)

// randomBigInt is a shortcut for generating a random big.Int
func randomBigInt() *big.Int {
	//Max random value, a 130-bits integer, i.e 2^130 - 1
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	//Generate cryptographically strong pseudo-random between 0 - max
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		//error handling
		return nil
	}

	return n
}

func TestTxSortedList(t *testing.T) {
	el := newTxSortedList()
	nItems := 100

	for i := 0; i < nItems; i++ {
		el.add(&TxTracker{HashStr: fmt.Sprintf("0x%d", i), GasPrice: randomBigInt()})
	}

	for i := 0; i < nItems-1; i++ {
		if !(el.getByIndex(i).GasPrice.Cmp(el.getByIndex(i+1).GasPrice) == 1) {
			t.Fatalf("Sort error. [%d].GasPrice(%f) < [%d].GasPrice(%f)", i, el.getByIndex(i).GasPrice, i+1, el.getByIndex(i+1).GasPrice)
		}
	}

	// el.print()

	if el.len() != nItems {
		t.Fatalf("Length error. Length %d. Expected %d", el.len(), nItems)
	}
}

func TestTxSortedListDelete(t *testing.T) {
	el := newTxSortedList()

	el.add(&TxTracker{HashStr: "0x01", GasPrice: new(big.Int).SetInt64(10)})
	el.add(&TxTracker{HashStr: "0x02", GasPrice: new(big.Int).SetInt64(20)})
	el.add(&TxTracker{HashStr: "0x03", GasPrice: new(big.Int).SetInt64(20)})
	el.add(&TxTracker{HashStr: "0x04", GasPrice: new(big.Int).SetInt64(40)})
	el.add(&TxTracker{HashStr: "0x05", GasPrice: new(big.Int).SetInt64(100)})
	el.add(&TxTracker{HashStr: "0x06", GasPrice: new(big.Int).SetInt64(15)})
	el.add(&TxTracker{HashStr: "0x07", GasPrice: new(big.Int).SetInt64(15)})
	el.add(&TxTracker{HashStr: "0x08", GasPrice: new(big.Int).SetInt64(10)})

	sort := []string{"0x05", "0x04", "0x02", "0x03", "0x06", "0x07", "0x01", "0x08"}

	for index, tx := range el.sorted {
		if sort[index] != tx.HashStr {
			t.Fatalf("Sort error. Expected %s, Actual %s", sort[index], tx.HashStr)
		}
	}

	deltxs := []string{"0x03", "0x06", "0x08", "0x05"}

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

	if el.delete(&TxTracker{HashStr: "0x08"}) {
		t.Fatal("Delete error. 0x08 tx was deleted and should not exist in the list")
	}
}

func TestTxSortedListBench(t *testing.T) {
	el := newTxSortedList()

	start := time.Now()
	for i := 0; i < 10000; i++ {
		el.add(&TxTracker{HashStr: fmt.Sprintf("0x%d", i), GasPrice: randomBigInt()})
	}
	elapsed := time.Since(start)
	t.Logf("TxSortedList adding 10000 items took %s", elapsed)

	start = time.Now()
	el.add(&TxTracker{HashStr: fmt.Sprintf("0x%d", 10001), GasPrice: randomBigInt()})
	elapsed = time.Since(start)
	t.Logf("TxSortedList adding the 10001 item (GasPrice=random) took %s", elapsed)

	start = time.Now()
	el.add(&TxTracker{HashStr: fmt.Sprintf("0x%d", 10002), GasPrice: new(big.Int).SetInt64(0)})
	elapsed = time.Since(start)
	t.Logf("TxSortedList adding the 10002 item (GasPrice=0) took %s", elapsed)

	start = time.Now()
	el.add(&TxTracker{HashStr: fmt.Sprintf("0x%d", 10003), GasPrice: new(big.Int).SetInt64(1000)})
	elapsed = time.Since(start)
	t.Logf("TxSortedList adding the 10003 item (GasPrice=1000) took %s", elapsed)
}
