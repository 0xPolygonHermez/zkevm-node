package broadcast_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast/mocks"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	host = "0.0.0.0"
	port = 61091
)

var (
	address      = fmt.Sprintf("%s:%d", host, port)
	broadcastSrv *broadcast.Server
	conn         *grpc.ClientConn
	cancel       context.CancelFunc
	err          error
	ctx          = context.Background()
)

func init() {
	// Change dir to project root
	// This is important because we have relative paths to files containing test vectors
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	initialize()
	defer teardown()

	os.Exit(m.Run())
}

func initialize() {
	broadcastSrv = initBroadcastServer()
	go broadcastSrv.Start()

	conn, cancel, err = initConn()
	if err != nil {
		panic(err)
	}

	err = operations.WaitGRPCHealthy(address)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	cancel()
	broadcastSrv.Stop()
}

func initConn() (*grpc.ClientConn, context.CancelFunc, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	conn, err := grpc.DialContext(ctx, address, opts...)
	return conn, cancel, err
}

func initBroadcastServer() *broadcast.Server {
	s := grpc.NewServer()
	st := new(mocks.StateMock)
	cfg := &broadcast.ServerConfig{
		Host: host,
		Port: port,
	}

	broadcastSrv = broadcast.NewServer(cfg, st)
	pb.RegisterBroadcastServiceServer(s, broadcastSrv)

	return broadcastSrv
}

func TestBroadcastServerGetBatch(t *testing.T) {
	tcs := []struct {
		description         string
		inputBatchNumber    uint64
		expectedBatch       *state.Batch
		expectedForcedBatch *state.ForcedBatch
		expectedEncodedTxs  []string
		expectedGER         *state.GlobalExitRoot
		expectedErr         bool
		expectedErrMsg      string
	}{
		{
			description:      "happy path",
			inputBatchNumber: 14,
			expectedBatch: &state.Batch{
				BatchNumber:    14,
				GlobalExitRoot: common.HexToHash("a"),
				Timestamp:      time.Now(),
			},
			expectedForcedBatch: &state.ForcedBatch{
				ForcedBatchNumber: 1,
			},
			expectedEncodedTxs: []string{"tx1", "tx2", "tx3"},
			expectedGER: &state.GlobalExitRoot{
				MainnetExitRoot: common.HexToHash("b"),
				RollupExitRoot:  common.HexToHash("c"),
			},
		},
		{
			description:      "query errors are returned",
			inputBatchNumber: 14,
			expectedErr:      true,
			expectedErrMsg:   "query error",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			st := new(mocks.StateMock)
			var err error
			if tc.expectedErr {
				err = errors.New(tc.expectedErrMsg)
			}
			st.On("GetBatchByNumber", mock.AnythingOfType("*context.valueCtx"), tc.inputBatchNumber, nil).Return(tc.expectedBatch, err)
			st.On("GetEncodedTransactionsByBatchNumber", mock.AnythingOfType("*context.valueCtx"), tc.inputBatchNumber, nil).Return(tc.expectedEncodedTxs, err)
			st.On("GetForcedBatchByBatchNumber", mock.AnythingOfType("*context.valueCtx"), tc.inputBatchNumber, nil).Return(tc.expectedForcedBatch, err)
			if tc.expectedBatch != nil {
				st.On("GetExitRootByGlobalExitRoot", mock.AnythingOfType("*context.valueCtx"), tc.expectedBatch.GlobalExitRoot, nil).Return(tc.expectedGER, err)
			}
			broadcastSrv.SetState(st)

			client := pb.NewBroadcastServiceClient(conn)
			actualBatch, err := client.GetBatch(ctx, &pb.GetBatchRequest{
				BatchNumber: tc.inputBatchNumber,
			})
			require.NoError(t, testutils.CheckError(err, tc.expectedErr, fmt.Sprintf("rpc error: code = Unknown desc = %s", tc.expectedErrMsg)))

			if err == nil {
				require.Equal(t, tc.expectedBatch.BatchNumber, actualBatch.BatchNumber)
				require.Equal(t, tc.expectedBatch.GlobalExitRoot.String(), actualBatch.GlobalExitRoot)
				require.Equal(t, uint64(tc.expectedBatch.Timestamp.Unix()), actualBatch.Timestamp)
				for i, encoded := range tc.expectedEncodedTxs {
					require.Equal(t, encoded, actualBatch.Transactions[i].Encoded)
				}
				require.Equal(t, tc.expectedGER.MainnetExitRoot.String(), actualBatch.MainnetExitRoot)
				require.Equal(t, tc.expectedGER.RollupExitRoot.String(), actualBatch.RollupExitRoot)

				require.True(t, st.AssertExpectations(t))
			}
		})
	}
}

func TestBroadcastServerGetLastBatch(t *testing.T) {
	tcs := []struct {
		description         string
		expectedBatch       *state.Batch
		expectedForcedBatch *state.ForcedBatch
		expectedEncodedTxs  []string
		expectedGER         *state.GlobalExitRoot
		expectedErr         bool
		expectedErrMsg      string
	}{
		{
			description: "happy path",
			expectedBatch: &state.Batch{
				BatchNumber:    14,
				GlobalExitRoot: common.HexToHash("a"),
				Timestamp:      time.Now(),
			},
			expectedForcedBatch: &state.ForcedBatch{
				ForcedBatchNumber: 1,
			},
			expectedEncodedTxs: []string{"tx1", "tx2", "tx3"},
			expectedGER: &state.GlobalExitRoot{
				MainnetExitRoot: common.HexToHash("b"),
				RollupExitRoot:  common.HexToHash("c"),
			},
		},
		{
			description:    "query errors are returned",
			expectedErr:    true,
			expectedErrMsg: "query error",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			st := new(mocks.StateMock)
			var err error
			if tc.expectedErr {
				err = errors.New(tc.expectedErrMsg)
			}
			st.On("GetLastBatch", mock.AnythingOfType("*context.valueCtx"), nil).Return(tc.expectedBatch, err)
			if tc.expectedBatch != nil {
				st.On("GetEncodedTransactionsByBatchNumber", mock.AnythingOfType("*context.valueCtx"), tc.expectedBatch.BatchNumber, nil).Return(tc.expectedEncodedTxs, err)
				st.On("GetForcedBatchByBatchNumber", mock.AnythingOfType("*context.valueCtx"), tc.expectedBatch.BatchNumber, nil).Return(tc.expectedForcedBatch, err)
				st.On("GetExitRootByGlobalExitRoot", mock.AnythingOfType("*context.valueCtx"), tc.expectedBatch.GlobalExitRoot, nil).Return(tc.expectedGER, err)
			}

			broadcastSrv.SetState(st)

			client := pb.NewBroadcastServiceClient(conn)
			actualBatch, err := client.GetLastBatch(ctx, &emptypb.Empty{})
			require.NoError(t, testutils.CheckError(err, tc.expectedErr, fmt.Sprintf("rpc error: code = Unknown desc = %s", tc.expectedErrMsg)))

			if err == nil {
				require.Equal(t, tc.expectedBatch.BatchNumber, actualBatch.BatchNumber)
				require.Equal(t, tc.expectedBatch.GlobalExitRoot.String(), actualBatch.GlobalExitRoot)
				require.Equal(t, uint64(tc.expectedBatch.Timestamp.Unix()), actualBatch.Timestamp)
				for i, encoded := range tc.expectedEncodedTxs {
					require.Equal(t, encoded, actualBatch.Transactions[i].Encoded)
				}
				require.Equal(t, tc.expectedGER.MainnetExitRoot.String(), actualBatch.MainnetExitRoot)
				require.Equal(t, tc.expectedGER.RollupExitRoot.String(), actualBatch.RollupExitRoot)

				require.True(t, st.AssertExpectations(t))
			}
		})
	}
}
