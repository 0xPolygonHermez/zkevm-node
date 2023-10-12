package gasprice

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateRate(t *testing.T) {
	testcases := []struct {
		l2CoinId int
		msg      string
		check    func(rate float64, err error)
	}{
		{
			// error
			l2CoinId: okbcoinId,
			msg:      "{\"topic\":\"middle_coinPrice_push\"}",
			check: func(rate float64, err error) {
				require.Error(t, err)
			},
		},
		{
			// error
			l2CoinId: okbcoinId,
			msg:      fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.02}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", okbcoinId+1),
			check: func(rate float64, err error) {
				require.Error(t, err)
			},
		},
		{
			// correct
			l2CoinId: okbcoinId,
			msg:      fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.02}, {\"coinId\":%d,\"price\":0.002}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", 1, okbcoinId),
			check: func(rate float64, err error) {
				require.Equal(t, rate, 0.002)
				require.NoError(t, err)
			},
		},
		{
			// correct
			l2CoinId: okbcoinId,
			msg:      fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.02}, {\"coinId\":%d,\"price\":10}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", 1, okbcoinId),
			check: func(rate float64, err error) {
				require.Equal(t, rate, float64(10))
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testcases {
		rp := newKafkaProcessor(Config{Topic: "middle_coinPrice_push", L2CoinId: tc.l2CoinId}, context.Background())
		rt, err := rp.parseL2CoinPrice([]byte(tc.msg))
		tc.check(rt, err)
	}
}
