package synchronizer

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	externalControlFilename = "/tmp/synchronizer_in"
	filePermissions         = 0644
	sleepTimeToReadFile     = 500 * time.Millisecond
)

// This is a local end-point in filesystem to send commands to a running synchronizer
// this is used for debugging purposes, to provide a way to reproduce some situations that are difficult
// to reproduce in a real test.
// It accept next commands:
// l1_producer_stop: stop producer
// l1_orchestrator_reset: reset orchestrator to a given block number
//
// example of usage (first you need to run the service):
// echo "l1_producer_stop" >> /tmp/synchronizer_in
// echo "l1_orchestrator_reset|8577060" >> /tmp/synchronizer_in
type externalControl struct {
	producer     *l1RollupInfoProducer
	orquestrator *l1SyncOrchestration
}

func newExternalControl(producer *l1RollupInfoProducer, orquestrator *l1SyncOrchestration) *externalControl {
	return &externalControl{producer: producer, orquestrator: orquestrator}
}

func (e *externalControl) start() {
	log.Infof("EXT:start: starting external control opening %s", externalControlFilename)
	file, err := os.OpenFile(externalControlFilename, os.O_APPEND|os.O_CREATE|os.O_RDONLY, filePermissions)
	if err != nil {
		log.Warnf("EXT:start:error opening file %s: %v", externalControlFilename, err)
		return
	}
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		log.Warnf("EXT:start:error seeking file %s: %v", externalControlFilename, err)
	}
	go e.readFile(file)
}

// https://medium.com/@arunprabhu.1/tailing-a-file-in-golang-72944204f22b
func (e *externalControl) readFile(file *os.File) {
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		for {
			line, err := reader.ReadString('\n')

			if err != nil {
				if err == io.EOF {
					// without this sleep you would hogg the CPU
					time.Sleep(sleepTimeToReadFile)
					continue
				}

				break
			}
			log.Infof("EXT:readFile: new command: %s", line)
			e.process(line)
		}
	}
}

func (e *externalControl) process(line string) {
	cmd := strings.Split(line, "|")
	if len(cmd) < 1 {
		return
	}
	switch strings.TrimSpace(cmd[0]) {
	case "l1_producer_stop":
		e.cmdL1ProducerStop(cmd[1:])
	case "l1_orchestrator_reset":
		e.cmdL1OrchestratorReset(cmd[1:])
	case "l1_orchestrator_stop":
		e.cmdL1OrchestratorAbort(cmd[1:])
	default:
		log.Warnf("EXT:process: unknown command: %s", cmd[0])
	}
}

func (e *externalControl) cmdL1OrchestratorReset(args []string) {
	log.Infof("EXT:cmdL1OrchestratorReset: %s", args)
	if len(args) < 1 {
		log.Infof("EXT:cmdL1OrchestratorReset: missing block number")
		return
	}
	blockNumber, err := strconv.ParseUint(strings.TrimSpace(args[0]), 10, 64)
	if err != nil {
		log.Infof("EXT:cmdL1OrchestratorReset: error parsing block number: %s", err)
		return
	}
	log.Infof("EXT:cmdL1OrchestratorReset: calling orchestrator reset(%d)", blockNumber)
	e.orquestrator.reset(blockNumber)
	log.Infof("EXT:cmdL1OrchestratorReset: calling orchestrator reset(%d) returned", blockNumber)
}

func (e *externalControl) cmdL1OrchestratorAbort(args []string) {
	log.Infof("EXT:cmdL1OrchestratorAbort: %s", args)
	if e.orquestrator == nil {
		log.Infof("EXT:cmdL1OrchestratorAbort: orquestrator is nil")
		return
	}
	log.Infof("EXT:cmdL1OrchestratorAbort: calling orquestrator stop")
	e.orquestrator.abort()
	log.Infof("EXT:cmdL1OrchestratorAbort: calling orquestrator stop returned")
}

func (e *externalControl) cmdL1ProducerStop(args []string) {
	log.Infof("EXT:cmdL1Stop: %s", args)
	if e.producer == nil {
		log.Infof("EXT:cmdL1Stop: producer is nil")
		return
	}
	log.Infof("EXT:cmdL1Stop: calling producer stop")
	e.producer.Stop()
	log.Infof("EXT:cmdL1Stop: calling producer stop returned")
}
