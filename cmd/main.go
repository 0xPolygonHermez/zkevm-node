package main

import (
	"os"
	"os/signal"

	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/mocks"
	"github.com/hermeznetwork/hermez-core/pool"
)

func main() {
	c := config.Load()
	setupLog(c.Log)
	runMigrations(c.Database)
	go runJSONRpcServer(c.RPC, c.Database)
	// go runSequencer(c.Sequencer)
	// go runAggregator(c.Aggregator)
	waitSignal()
}

func setupLog(c log.Config) {
	log.Init(c)
}

func runMigrations(c db.Config) {
	err := db.RunMigrations(c)
	if err != nil {
		log.Fatal(err)
	}
}

func runJSONRpcServer(jc jsonrpc.Config, dc db.Config) {
	p, err := pool.NewPostgresPool(dc)
	if err != nil {
		log.Fatal(err)
	}

	s := mocks.NewState()

	if err := jsonrpc.NewServer(jc, p, s).Start(); err != nil {
		log.Fatal(err)
	}
}

// func runSequencer(c sequencer.Config) {
// 	p := mocks.NewPool()
// 	s := mocks.NewState()
// 	e, err := etherman.NewEtherman(c.Etherman)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	sy, err := synchronizer.NewSynchronizer(e, s)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	seq, err := sequencer.NewSequencer(c, p, s, e, sy)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	seq.Start()
// }

// func runAggregator(c aggregator.Config) {
// 	// TODO: have more readable variables
// 	s := mocks.NewState()
// 	bp := s.NewBatchProcessor(common.Hash{}, false)
// 	e, err := etherman.NewEtherman(c.Etherman)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	sy, err := synchronizer.NewSynchronizer(e, s)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	pc := aggregator.NewProverClient()
// 	agg, err := aggregator.NewAggregator(c, s, bp, e, sy, pc)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	agg.Start()
// }

func waitSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for sig := range signals {
		switch sig {
		case os.Interrupt, os.Kill:
			log.Info("terminating application gracefully...")
			os.Exit(0)
		}
	}
}
