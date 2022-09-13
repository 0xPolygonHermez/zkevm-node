package operations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// DefaultInterval is a time interval
	DefaultInterval = 2 * time.Millisecond
	// DefaultDeadline is a time interval
	DefaultDeadline = 30 * time.Second
	// DefaultTxMinedDeadline is a time interval
	DefaultTxMinedDeadline = 5 * time.Second
)

var (
	// ErrTimeoutReached is thrown when the timeout is reached and
	// because the condition is not matched
	ErrTimeoutReached = fmt.Errorf("timeout has been reached")
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
			return ErrTimeoutReached
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

type ethClienter interface {
	ethereum.TransactionReader
	ethereum.ContractCaller
}

// WaitTxToBeMined waits until a tx has been mined or the given timeout expires.
func WaitTxToBeMined(client ethClienter, hash common.Hash, timeout time.Duration) error {
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

// WaitL2BlockToBeConsolidated waits until a L2 Block has been consolidated or the given timeout expires.
func WaitL2BlockToBeConsolidated(l2Block *big.Int, timeout time.Duration) error {
	return Poll(DefaultInterval, timeout, func() (bool, error) {
		return l2BlockConsolidationCondition(l2Block)
	})
}

// WaitL2BlockToBeVirtualized waits until a L2 Block has been virtualized or the given timeout expires.
func WaitL2BlockToBeVirtualized(l2Block *big.Int, timeout time.Duration) error {
	return Poll(DefaultInterval, timeout, func() (bool, error) {
		return l2BlockVirtualizationCondition(l2Block)
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
	return NodeUpCondition(DefaultL1NetworkURL)
}

func nodeUpCondition() (done bool, err error) {
	return NodeUpCondition(DefaultL2NetworkURL)
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
func txMinedCondition(ctx context.Context, client ethClienter, hash common.Hash) (bool, error) {
	// Get tx status
	tx, isPending, err := client.TransactionByHash(ctx, hash)
	if err == ethereum.NotFound || isPending {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	// Check if tx has failed
	receipt, err := client.TransactionReceipt(ctx, hash)
	if err != nil {
		return false, err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		// Get revert reason
		reason, reasonErr := revertReason(ctx, client, tx, receipt.BlockNumber)
		if reasonErr != nil {
			reason = reasonErr.Error()
		}
		return false, fmt.Errorf("transaction has failed, reason: %s, receipt: %+v. tx: %+v, gas: %v", reason, receipt, tx, tx.Gas())
	}
	return true, nil
}

func revertReason(ctx context.Context, c ethClienter, tx *types.Transaction, blockNumber *big.Int) (string, error) {
	from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		signer := types.LatestSignerForChainID(tx.ChainId())
		from, err = types.Sender(signer, tx)
		if err != nil {
			return "", err
		}
	}
	msg := ethereum.CallMsg{
		From: from,
		To:   tx.To(),
		Gas:  tx.Gas(),

		Value: tx.Value(),
		Data:  tx.Data(),
	}
	hex, err := c.CallContract(ctx, msg, blockNumber)
	if err != nil {
		return "", err
	}

	reasonOffset := new(big.Int).SetBytes(hex[4 : 4+32])
	reason := string(hex[4+32+int(reasonOffset.Uint64()):])
	return reason, nil
}

// l2BlockConsolidationCondition
func l2BlockConsolidationCondition(l2Block *big.Int) (bool, error) {
	l2NetworkURL := "http://localhost:8123"
	response, err := jsonrpc.JSONRPCCall(l2NetworkURL, "zkevm_isL2BlockConsolidated", l2Block.Uint64())
	if err != nil {
		return false, err
	}
	if response.Error != nil {
		return false, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}
	var result bool
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

// l2BlockVirtualizationCondition
func l2BlockVirtualizationCondition(l2Block *big.Int) (bool, error) {
	l2NetworkURL := "http://localhost:8123"
	response, err := jsonrpc.JSONRPCCall(l2NetworkURL, "zkevm_isL2BlockVirtualized", l2Block.Uint64())
	if err != nil {
		return false, err
	}
	if response.Error != nil {
		return false, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}
	var result bool
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return false, err
	}
	return result, nil
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
