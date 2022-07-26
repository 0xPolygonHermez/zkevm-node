package operations

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// DefaultInterval is a time interval
	DefaultInterval = 2 * time.Second
	// DefaultDeadline is a time interval
	DefaultDeadline = 30 * time.Second
	// DefaultTxMinedDeadline is a time interval
	DefaultTxMinedDeadline = 5 * time.Second
)

// Wait handles polliing until conditions are met.
type Wait struct{}

// NewWait is the Wait constructor.
func NewWait() *Wait {
	return &Wait{}
}

// Poll retries the given condition with the given interval until it succeeds
// or the given deadline expires.
func Poll(interval, deadline time.Duration, condition ConditionFunc) error {
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

// WaitTxToBeMined waits until a tx has been mined or the given timeout expires.
func WaitTxToBeMined(client *ethclient.Client, hash common.Hash, timeout time.Duration) error {
	ctx := context.Background()
	return Poll(DefaultInterval, timeout, func() (bool, error) {
		return txMinedCondition(ctx, client, hash)
	})
}

// WaitGRPCHealthy waits for a gRPC endpoint to be responding according to the
// health standard in package grpc.health.v1
func WaitGRPCHealthy(address string) error {
	return Poll(DefaultInterval, DefaultDeadline, func() (bool, error) {
		return grpcHealthyCondition(address)
	})
}

// NodeUpCondition check if the container is up and running
func NodeUpCondition(target string) (bool, error) {
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

// ConditionFunc is a generic function
type ConditionFunc func() (done bool, err error)

func networkUpCondition() (bool, error) {
	return NodeUpCondition(l1NetworkURL)
}

// ProverUpCondition check if the prover is up and running
func ProverUpCondition() (bool, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

	proverClient := pb.NewZKProverServiceClient(conn)
	state, err := proverClient.GetStatus(context.Background(), &pb.GetStatusRequest{})
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}

	done := state.State == pb.GetStatusResponse_STATUS_PROVER_IDLE

	return done, nil
}

func nodeUpCondition() (done bool, err error) {
	return NodeUpCondition(l2NetworkURL)
}

func grpcHealthyCondition(address string) (bool, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

// txMinedCondition
func txMinedCondition(ctx context.Context, client *ethclient.Client, hash common.Hash) (bool, error) {
	_, isPending, err := client.TransactionByHash(ctx, hash)
	if err == ethereum.NotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	var done bool
	if !isPending {
		r, err := client.TransactionReceipt(ctx, hash)
		if err != nil {
			return false, err
		}
		if r.Status == types.ReceiptStatusFailed {
			return false, fmt.Errorf("transaction has failed: %s", hex.EncodeToString(r.PostState))
		}
		done = true
	}
	return done, nil
}

// WaitSignal blocks until an Interrupt or Kill signal is received, then it
// executes the given cleanup functions and returns.
func WaitSignal(cleanupFuncs ...func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for sig := range signals {
		switch sig {
		case os.Interrupt, os.Kill:
			log.Info("terminating application gracefully...")
			for _, cleanup := range cleanupFuncs {
				cleanup()
			}
			return
		}
	}
}
