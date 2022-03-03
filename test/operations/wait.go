package operations

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	defaultInterval        = 2 * time.Second
	defaultDeadline        = 30 * time.Second
	defaultTxMinedDeadline = 5 * time.Second
)

// Wait handles polliing until conditions are met.
type Wait struct{}

// NewWait is the Wait constructor.
func NewWait() *Wait {
	return &Wait{}
}

// Poll retries the given condition with the given interval until it succeeds
// or the given deadline expires.
func (w *Wait) Poll(interval, deadline time.Duration, condition conditionFunc) error {
	timeout := time.After(deadline)
	tick := time.NewTicker(interval)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("Condition not met after %s", deadline)
		case <-tick.C:
			ok, err := condition()
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
		}
	}
}

// GRPCHealthy waits for a gRPC endpoint to be responding according to the
// health standard in package grpc.health.v1
func (w *Wait) GRPCHealthy(address string) error {
	return w.Poll(defaultInterval, defaultDeadline, func() (bool, error) {
		return grpcHealthyCondition(address)
	})
}

// TxToBeMined waits until a tx has been mined or the given timeout expires.
func (w *Wait) TxToBeMined(client *ethclient.Client, hash common.Hash, timeout time.Duration) error {
	start := time.Now()
	ctx := context.Background()
	for {
		if time.Since(start) > timeout {
			return errors.New("timeout exceed")
		}

		time.Sleep(1 * time.Second)

		_, isPending, err := client.TransactionByHash(ctx, hash)
		if err == ethereum.NotFound {
			continue
		}

		if err != nil {
			return err
		}

		if !isPending {
			r, err := client.TransactionReceipt(ctx, hash)
			if err != nil {
				return err
			}

			if r.Status == types.ReceiptStatusFailed {
				return fmt.Errorf("transaction has failed: %s", string(r.PostState))
			}

			return nil
		}
	}
}

// WaitGRPCHealthy waits for a gRPC endpoint to be responding according to the
// health standard in package grpc.health.v1
func WaitGRPCHealthy(address string) error {
	w := NewWait()
	return w.Poll(defaultInterval, defaultDeadline, func() (bool, error) {
		return grpcHealthyCondition(address)
	})
}

func nodeUpCondition(target string) (bool, error) {
	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}`)
	req, err := http.NewRequest(
		"POST", target,
		bytes.NewBuffer(jsonStr))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}

	if res.Body != nil {
		defer func() {
			err = res.Body.Close()
		}()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return false, err
	}

	r := struct {
		Result bool
	}{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return false, err
	}

	done := !r.Result

	return done, nil
}

type conditionFunc func() (done bool, err error)

func networkUpCondition() (bool, error) {
	return nodeUpCondition(l1NetworkURL)
}

func proverUpCondition() (bool, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "localhost:50051", opts...)
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}
	defer func() {
		err = conn.Close()
	}()

	proverClient := proverclient.NewZKProverClient(conn)
	state, err := proverClient.GetStatus(context.Background(), &proverclient.NoParams{})
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}

	done := state.Status == proverclient.State_IDLE

	return done, nil
}

func coreUpCondition() (done bool, err error) {
	return nodeUpCondition(l2NetworkURL)
}

func grpcHealthyCondition(address string) (bool, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}
	defer func() {
		err = conn.Close()
	}()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	state, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}

	done := state.Status == grpc_health_v1.HealthCheckResponse_SERVING

	return done, nil
}
