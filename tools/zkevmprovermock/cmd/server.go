package main

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
	stateDBPpb "github.com/0xPolygonHermez/zkevm-node/merkletree/pb"
	executorPb "github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/tools/zkevmprovermock/server"
	"github.com/0xPolygonHermez/zkevm-node/tools/zkevmprovermock/testvector"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func runServer(cliCtx *cli.Context) error {
	log.Info("Running zkEVM Prover Server mock...")

	s := grpc.NewServer()

	aferoFs := afero.NewOsFs()

	tvContainer, err := testvector.NewContainer(cliCtx.String("test-vector-path"), aferoFs)
	if err != nil {
		log.Fatalf("Could not create test vector container: %v", err)
	}

	stateDBAddress := fmt.Sprintf("%s:%d", cliCtx.String("host"), cliCtx.Uint("statedb-port"))
	stateDBSrv := server.NewStateDBMock(stateDBAddress, tvContainer)

	stateDBPpb.RegisterStateDBServiceServer(s, stateDBSrv)
	go stateDBSrv.Start()

	executorAddress := fmt.Sprintf("%s:%d", cliCtx.String("host"), cliCtx.Uint("executor-port"))
	executorSrv := server.NewExecutorMock(executorAddress, tvContainer)

	executorPb.RegisterExecutorServiceServer(s, executorSrv)
	go executorSrv.Start()

	operations.WaitSignal(func() {
		stateDBSrv.Stop()
		executorSrv.Stop()
	})
	return nil
}
